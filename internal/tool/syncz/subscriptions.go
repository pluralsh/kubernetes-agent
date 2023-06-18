package syncz

import (
	"context"
	"sync"
)

type EventCallback[E any] func(ctx context.Context, e E)

type Subscriptions[E any] struct {
	mu  sync.Mutex
	chs []chan<- E
}

func (s *Subscriptions[E]) add(ch chan<- E) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chs = append(s.chs, ch)
}

func (s *Subscriptions[E]) remove(ch chan<- E) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, c := range s.chs {
		if c == ch {
			l := len(s.chs)
			newChs := append(s.chs[:i], s.chs[i+1:]...)
			s.chs[l-1] = nil // help GC
			s.chs = newChs
			break
		}
	}
}

func (s *Subscriptions[E]) On(ctx context.Context, cb EventCallback[E]) {
	ch := make(chan E)
	s.add(ch)
	defer s.remove(ch)

	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		case e := <-ch:
			cb(ctx, e)
		}
	}
}

// Dispatch dispatches the given event to all added subscriptions.
func (s *Subscriptions[E]) Dispatch(ctx context.Context, e E) {
	done := ctx.Done()

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.chs {
		select {
		case <-done:
			return
		case ch <- e:
		}
	}
}
