package agent

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/wait"
)

type WorkSource interface {
	ID() string
	Configuration() proto.Message
}

type WorkerFactory interface {
	New(agentId int64, source WorkSource) Worker
	SourcesFromConfiguration(*agentcfg.AgentConfiguration) []WorkSource
}

type Worker interface {
	Run(context.Context)
}

type WorkerManager struct {
	log           *zap.Logger
	workerFactory WorkerFactory
	workers       map[string]*workerHolder // source id -> worker holder instance
}

func NewWorkerManager(log *zap.Logger, workerFactory WorkerFactory) *WorkerManager {
	return &WorkerManager{
		log:           log,
		workerFactory: workerFactory,
		workers:       map[string]*workerHolder{},
	}
}

func (m *WorkerManager) startNewWorker(agentId int64, source WorkSource) {
	id := source.ID()
	m.log.Info("Starting synchronization worker", logz.WorkerId(id))
	worker := m.workerFactory.New(agentId, source)
	ctx, cancel := context.WithCancel(context.Background())
	holder := &workerHolder{
		source: source,
		stop:   cancel,
	}
	holder.wg.StartWithContext(ctx, worker.Run)
	m.workers[id] = holder
}

func (m *WorkerManager) ApplyConfiguration(agentId int64, cfg *agentcfg.AgentConfiguration) error {
	sources := m.workerFactory.SourcesFromConfiguration(cfg)
	newSetOfSources := make(map[string]struct{}, len(sources))
	var sourcesToStartWorkersFor []WorkSource
	var workersToStop []*workerHolder //nolint:prealloc

	// Collect sources without workers or with updated configuration.
	for _, source := range sources {
		id := source.ID()
		if _, ok := newSetOfSources[id]; ok {
			return fmt.Errorf("duplicate source id: %s", id)
		}
		newSetOfSources[id] = struct{}{}
		holder := m.workers[id]
		if holder == nil { // New source added
			sourcesToStartWorkersFor = append(sourcesToStartWorkersFor, source)
		} else { // We have a worker for this source already
			if proto.Equal(source.Configuration(), holder.source.Configuration()) {
				// Worker's configuration hasn't changed, nothing to do here
				continue
			}
			m.log.Info("Configuration has been updated, restarting synchronization worker", logz.WorkerId(id))
			workersToStop = append(workersToStop, holder)
			sourcesToStartWorkersFor = append(sourcesToStartWorkersFor, source)
		}
	}

	// Stop workers for sources which have been removed from the list.
	for sourceId, holder := range m.workers {
		if _, ok := newSetOfSources[sourceId]; ok {
			continue
		}
		workersToStop = append(workersToStop, holder)
	}

	// Tell workers that should be stopped to stop.
	for _, holder := range workersToStop {
		m.log.Info("Stopping synchronization worker", logz.WorkerId(holder.source.ID()))
		holder.stop()
		delete(m.workers, holder.source.ID())
	}

	// Wait for stopped workers to finish.
	for _, holder := range workersToStop {
		m.log.Info("Waiting for synchronization worker to stop", logz.WorkerId(holder.source.ID()))
		holder.wg.Wait()
	}

	// Start new workers for new sources or because of updated configuration.
	for _, source := range sourcesToStartWorkersFor {
		m.startNewWorker(agentId, source) // nolint: contextcheck
	}
	return nil
}

func (m *WorkerManager) StopAllWorkers() {
	// Tell all workers to stop
	for _, holder := range m.workers {
		holder.stop()
	}
	// Wait for all workers to stop
	for _, holder := range m.workers {
		holder.wg.Wait()
	}
}

type workerHolder struct {
	source WorkSource
	wg     wait.Group
	stop   context.CancelFunc
}
