package syncz

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestSubscriptions_DispatchingMultiple(t *testing.T) {
	// GIVEN
	var wg wait.Group
	defer wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var s Subscriptions[int]
	x := 42
	// recorder for callback hits
	rec1 := make(chan struct{})
	rec2 := make(chan struct{})
	subscriber1 := func(_ context.Context, e int) {
		assert.EqualValues(t, x, e)
		close(rec1)
	}
	subscriber2 := func(_ context.Context, e int) {
		assert.EqualValues(t, x, e)
		close(rec2)
	}

	// WHEN
	// starting multiple subscribers
	wg.Start(func() {
		s.On(ctx, subscriber1)
	})
	wg.Start(func() {
		s.On(ctx, subscriber2)
	})

	// give the OnGitPushEvent goroutines time to be scheduled and registered
	assert.Eventually(t, func() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		return len(s.chs) == 2
	}, time.Minute, time.Millisecond)

	// dispatch a single event
	s.Dispatch(ctx, x)

	// THEN
	<-rec1
	<-rec2
}

func TestSubscriptions_AddRemove(t *testing.T) {
	var s Subscriptions[int]

	ch1 := make(chan<- int)
	ch2 := make(chan<- int)
	ch3 := make(chan<- int)

	s.add(ch1)
	s.add(ch2)
	s.add(ch3)

	assert.Equal(t, ch1, s.chs[0])
	assert.Equal(t, ch2, s.chs[1])
	assert.Equal(t, ch3, s.chs[2])

	s.remove(ch2)

	assert.Equal(t, ch1, s.chs[0])
	assert.Equal(t, ch3, s.chs[1])
	assert.Nil(t, s.chs[:3:3][2])

	s.remove(ch1)
	s.remove(ch3)
	assert.Nil(t, s.chs[:3:3][0])
	assert.Nil(t, s.chs[:3:3][1])
	assert.Empty(t, s.chs)
}
