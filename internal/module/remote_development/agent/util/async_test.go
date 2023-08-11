package util

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunWithAsyncResult_SimplePublish(t *testing.T) {
	var wg sync.WaitGroup

	// ensure async func blocks during execution
	wg.Add(1)
	resultCh := RunWithAsyncResult(func(ch chan<- int) {
		ch <- 1

		wg.Wait()
		ch <- 2
	})

	expectChannelToReceive(t, resultCh, 1)
	expectChannelToBeBlocked(t, resultCh)

	// resume execution by unblocking wait-group
	wg.Done()

	expectChannelToReceive(t, resultCh, 2)
	expectChannelToBeClosed(t, resultCh)
}

func TestToAsync(t *testing.T) {
	value := 1
	asyncValue := ToAsync(value)

	expectChannelToReceive(t, asyncValue, value)
	expectChannelToBeClosed(t, asyncValue)
}

func TestCombineChannels_ClosedChannel(t *testing.T) {
	ch1 := RunWithAsyncResult(func(ch chan<- int) {
		ch <- 1
		ch <- 2
	})

	ch2 := make(chan int)
	close(ch2)

	channels := []<-chan int{ch1, ch2}
	resultCh := CombineChannels(channels)

	var received []int //nolint:prealloc
	for value := range resultCh {
		received = append(received, value)
	}

	require.Len(t, received, 2)
	require.Contains(t, received, 1)
	require.Contains(t, received, 2)
}

func TestCombineChannels_DelayedPublish(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	ch1 := RunWithAsyncResult(func(ch chan<- int) {
		wg.Wait()
		ch <- 1
		ch <- 2
	})

	ch2 := RunWithAsyncResult(func(ch chan<- int) {
		ch <- 3
	})

	channels := []<-chan int{ch1, ch2}
	resultCh := CombineChannels(channels)

	expectChannelToReceive(t, resultCh, 3)

	// unblock remaining publish
	wg.Done()

	expectChannelToReceive(t, resultCh, 1)
	expectChannelToReceive(t, resultCh, 2)
	expectChannelToBeClosed(t, resultCh)
}

func TestCombineChannels_NilSlice(t *testing.T) {
	resultCh := CombineChannels[int](nil)

	expectChannelToBeClosed(t, resultCh)
}

func expectChannelToReceive[T any](t *testing.T, ch <-chan T, expected T) {
	actual := <-ch
	require.Equal(t, expected, actual)
}

func expectChannelToBeBlocked[T any](t *testing.T, ch <-chan T) {
	select {
	case received := <-ch:
		t.Errorf("Value %v received but not expected", received)
	default:
		// the default block is expected to execute as the channel is blocked
	}
}

func expectChannelToBeClosed[T any](t *testing.T, ch <-chan T) {
	_, isOpen := <-ch
	require.False(t, isOpen)
}
