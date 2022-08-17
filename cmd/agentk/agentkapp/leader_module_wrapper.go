package agentkapp

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
)

type Runner interface {
	// RunWhenLeader runs f when it is elected as the leader.
	// Provided context signals done when leadership is lost or when the program is terminating.
	// f may be called, stopped, called multiple times, depending on leadership status.
	// The returned function can be used to signal f to stop if it's already running or to avoid running it
	// if elected as the leader. The returned function blocks until f returns.
	RunWhenLeader(f func(context.Context)) func()
}

type leaderModuleWrapper struct {
	module        modagent.LeaderModule
	runner        Runner
	cfg2module    chan *agentcfg.AgentConfiguration
	errFromModule chan error
	stopRun       func() // used to stop the maybe running module and as a flag.
}

func newLeaderModuleWrapper(module modagent.LeaderModule, runner Runner) *leaderModuleWrapper {
	return &leaderModuleWrapper{
		module: module,
		runner: runner,
	}
}

func (w *leaderModuleWrapper) DefaultAndValidateConfiguration(cfg *agentcfg.AgentConfiguration) error {
	return w.module.DefaultAndValidateConfiguration(cfg)
}

func (w *leaderModuleWrapper) Name() string {
	return w.module.Name()
}

func (w *leaderModuleWrapper) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) (retErr error) {
	var (
		nilableCfg chan<- *agentcfg.AgentConfiguration
		config     *agentcfg.AgentConfiguration
	)

	defer func() {
		err := w.stopWithErr()
		if retErr == nil {
			retErr = err
		}
	}()
	for {
		select {
		case c, ok := <-cfg: // case #1
			if !ok {
				return nil
			}
			if !w.module.IsRunnableConfiguration(c) {
				err := w.stopWithErr()
				if err != nil {
					return err
				}
				nilableCfg = nil // disable case #2
				continue
			}
			if w.stopRun == nil { // Not running yet
				w.run()
			}
			config = c
			nilableCfg = w.cfg2module // enable case #2
		case nilableCfg <- config: // case #2, disabled when nilableCfg == nil i.e. when there is nothing to send
			// config sent
			config = nil     // help GC
			nilableCfg = nil // disable case #2
		case err := <-w.errFromModule:
			// If we are here it means the module returned without us asking.
			w.stop() // Clean things up after early return
			if err != nil {
				return err
			}
			if config != nil {
				// We have a config that we haven't sent to the module while it was running.
				// Let's do that. Start the module and then the next loop will send the config.
				w.run()
				nilableCfg = w.cfg2module // enable case #2
			} else {
				w.errFromModule = nil // disable this select case if module is not running
				nilableCfg = nil      // disable case #2 if module is not running
			}
		}
	}
}

func (w *leaderModuleWrapper) run() {
	w.cfg2module = make(chan *agentcfg.AgentConfiguration)
	w.errFromModule = make(chan error, 1)
	w.stopRun = w.runner.RunWhenLeader(func(ctx context.Context) {
		w.errFromModule <- w.module.Run(ctx, w.cfg2module)
	})
}

func (w *leaderModuleWrapper) stop() bool {
	if w.stopRun == nil {
		return false
	}
	close(w.cfg2module)
	w.cfg2module = nil
	w.stopRun()
	w.stopRun = nil
	return true
}

func (w *leaderModuleWrapper) stopWithErr() error {
	if !w.stop() {
		return nil
	}
	close(w.errFromModule)
	// either we get a value or nil if there is no value (module was not running) and the channel was closed
	err := <-w.errFromModule
	w.errFromModule = nil
	return err
}
