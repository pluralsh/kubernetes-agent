package agentkapp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type electorStatus byte

const (
	notRunning electorStatus = iota
	runningButNotLeader
	runningAndLeader
	stopping
)

type funcHolder struct {
	f      func(context.Context)
	cancel context.CancelFunc
	wait   func()
}

func (h *funcHolder) start() {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel
	h.wait = wg.Wait
	go func() {
		defer wg.Done()
		defer cancel()
		h.f(ctx)
	}()
}

type LeaderElector interface {
	Run(ctx context.Context, onStartedLeading, onStoppedLeading func())
}

type leaderRunner struct {
	leaderElector    LeaderElector
	addFunc          chan addFuncReq
	stopFunc         chan stopFuncReq
	onStartedLeading chan struct{}
	onStoppedLeading chan struct{}
	funcs            map[int32]*funcHolder
	idxCounter       int32
	status           electorStatus
}

type addFuncReq struct {
	f        func(context.Context)
	stopResp chan<- func()
}

type stopFuncReq struct {
	idx         int32
	waitForStop chan<- func()
}

func newLeaderRunner(leaderElector LeaderElector) *leaderRunner {
	return &leaderRunner{
		leaderElector:    leaderElector,
		addFunc:          make(chan addFuncReq),
		stopFunc:         make(chan stopFuncReq),
		onStartedLeading: make(chan struct{}),
		onStoppedLeading: make(chan struct{}),
		funcs:            make(map[int32]*funcHolder),
		status:           notRunning,
	}
}

func (r *leaderRunner) MaybeWrapModule(m modagent.Module) modagent.Module {
	lm, ok := m.(modagent.LeaderModule)
	if !ok {
		return m
	}
	return newLeaderModuleWrapper(lm, r)
}

func (r *leaderRunner) RunWhenLeader(f func(context.Context)) func() {
	stopResp := make(chan func())
	r.addFunc <- addFuncReq{
		f:        f,
		stopResp: stopResp,
	}
	return <-stopResp
}

// Run starts the leader election process and starts and stops the modules depending on if it is the leader or not.
func (r *leaderRunner) Run(ctx context.Context) {
	if r.status != notRunning {
		panic(fmt.Errorf("the leader runner can only be started once; status is %d", r.status))
	}
	defer func() {
		r.status = notRunning
	}()

	r.status = runningButNotLeader

	var wg wait.Group
	defer wg.Wait()
	wg.Start(func() {
		r.leaderElector.Run(ctx,
			func() { r.onStartedLeading <- struct{}{} },
			func() { r.onStoppedLeading <- struct{}{} },
		)
	})

	done := ctx.Done()
	for {
		select {
		case <-done:
			if len(r.funcs) > 0 {
				// This can only happen if leaderRunner is misused.
				// All funcs must be stopped before the context can signal done.
				panic(fmt.Errorf("%d functions still want to run; status is %d", len(r.funcs), r.status))
			}

			switch r.status {
			case stopping:
				<-r.onStoppedLeading // wait for elector to fully stop
			case notRunning:
				fallthrough
			case runningButNotLeader:
				fallthrough
			case runningAndLeader:
				fallthrough
			default:
				panic(fmt.Errorf("unexpected status: %d. Expecting stopping (%d)", r.status, stopping))
			}

			return
		case add := <-r.addFunc:
			idx := r.idxCounter
			r.idxCounter++
			holder := &funcHolder{
				f:      add.f,
				cancel: func() {},
				wait:   func() {},
			}
			switch r.status {
			case runningAndLeader:
				holder.start() // nolint: contextcheck
			case runningButNotLeader:
				// Nothing to do right now
				// The holder function will be started when elected as leader and
				// leader runner transitions to runningAndLeader
			case stopping:
				// Nothing to do right now
				// The holder function will be started when elected as leader and
				// leader runner transitions to runningAndLeader
			case notRunning:
				fallthrough
			default:
				panic(fmt.Errorf("unexpected status: %d", r.status))
			}
			r.funcs[idx] = holder
			add.stopResp <- func() {
				w := make(chan func())
				r.stopFunc <- stopFuncReq{
					idx:         idx,
					waitForStop: w,
				}
				(<-w)()
			}
		case stop := <-r.stopFunc:
			var stopFunc func()
			holder, ok := r.funcs[stop.idx]
			if ok {
				holder.cancel()
				delete(r.funcs, stop.idx)
				stopFunc = holder.wait
				if len(r.funcs) == 0 { // removed last function
					r.status = stopping
				}
			} else { // stop() called more than once
				stopFunc = func() {}
			}
			stop.waitForStop <- stopFunc
		case <-r.onStartedLeading:
			// Leader elector calls onStartedLeading() callback asynchronously in a goroutine.
			// Because of that there is no ordering guarantee on when this is triggered. To mitigate the race
			// we have to check the current state and only act on the notification if it is expected.
			switch r.status {
			case runningButNotLeader:
				r.status = runningAndLeader
				for _, holder := range r.funcs {
					holder.start() // nolint: contextcheck
				}
			case runningAndLeader:
				// This can happen if elector is stopped and started again really quickly and callback from the second
				// election executes before the callback from the first one. This is almost impossible.
				// Nothing to do here.
			case stopping:
				// Elected as leader then quickly stopped. Callback executed after elector has stopped.
				// Nothing to do here.
			case notRunning:
				fallthrough
			default:
				panic(fmt.Errorf("unexpected status: %d", r.status))
			}
		case <-r.onStoppedLeading:
			switch r.status {
			case runningButNotLeader:
				// onStoppedLeading() is called even if wasn't the leader.
				// Nothing to do here.
			case runningAndLeader:
				// Lost election. Must stop all functions ASAP.
				for _, holder := range r.funcs {
					holder.cancel()
				}
				for idx, holder := range r.funcs {
					holder.cancel = func() {}
					holder.wait()
					holder.wait = func() {}
					r.funcs[idx] = holder
				}
			case stopping:
				// This happens when last function was stopped.
				// Nothing to do here.
			case notRunning:
				fallthrough
			default:
				panic(fmt.Errorf("unexpected status: %d", r.status))
			}

			// We are not leading anymore.
			r.status = runningButNotLeader
		}
	}
}

type leaseLeaderElector struct {
	// Namespace is the namespace of the Lease lock object.
	namespace string
	// name returns name of the Lease lock object or an error if context signals done.
	name func(context.Context) (string, error)

	// Identity is the unique string identifying a lease holder across
	// all participants in an election.
	identity string

	coordinationClient coordinationv1.CoordinationV1Interface
	eventRecorder      resourcelock.EventRecorder
}

func (l *leaseLeaderElector) Run(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
	name, err := l.name(ctx)
	if err != nil {
		return // ctx done
	}
	elector, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock: &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Namespace: l.namespace,
				Name:      name,
			},
			Client: l.coordinationClient,
			LockConfig: resourcelock.ResourceLockConfig{
				Identity:      l.identity,
				EventRecorder: l.eventRecorder,
			},
		},
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) { onStartedLeading() },
			OnStoppedLeading: onStoppedLeading,
		},
		ReleaseOnCancel: true,
		Name:            "module-runner",
	})
	if err != nil {
		// This can only happen if config is incorrect. It is hard-coded here, so should never happen.
		panic(err)
	}
	wait.UntilWithContext(ctx, elector.Run, time.Second)
}
