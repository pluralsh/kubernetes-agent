package util

import "k8s.io/apimachinery/pkg/util/wait"

// RunWithAsyncResult runs the passed function asynchronously and returns the results
// in a channel. The function may choose to publish as many values using the write-only channel in the
// callback. The write-only channel is closed automatically after the callback function finishes execution.
// The channel returned by the function is closed automatically after the last value has been read from the channel
func RunWithAsyncResult[T any](fn func(ch chan<- T)) <-chan T {
	responseCh := make(chan T)
	go func() {
		defer close(responseCh)

		fn(responseCh)
	}()
	return responseCh
}

// ToAsync takes any value and makes it available in a read-only channel. The channel is
// closed automatically after the value has been read from the channel
func ToAsync[T any](value T) <-chan T {
	c := make(chan T, 1)
	c <- value
	close(c)
	return c
}

// CombineChannels take a slice of non-nil channels as input and returns a chan that relays the
// values received across all the input channels. The values transmitted by the returned channel
// may be received out of order if they are transmitted by different input channels but will retain the order
// if they are transmitted by the same input channel. The returned channel is closed automatically
// after the last channel in the input channel list is closed
func CombineChannels[T any](channels []<-chan T) <-chan T {
	outCh := make(chan T)
	var wg wait.Group

	for _, ch := range channels {
		currentCh := ch
		wg.Start(func() {
			for value := range currentCh {
				outCh <- value
			}
		})
	}

	go func() {
		// only close the channel once each input channel is done publishing
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
