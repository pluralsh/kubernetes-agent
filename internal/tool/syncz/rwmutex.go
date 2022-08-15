package syncz

import "context"

type RWMutex struct {
	// See Mutex for the box analogy. This channel is use for grabbing exclusive access to the mutex.
	box chan struct{}
	// readers holds the current reader count. Empty channel means no readers OR that the number of readers
	// has been taken out to be modified.
	readers chan int32
}

func NewRWMutex() RWMutex {
	return RWMutex{
		box:     make(chan struct{}, 1),
		readers: make(chan int32, 1),
	}
}

func (m RWMutex) Lock(ctx context.Context) bool {
	select {
	case <-ctx.Done(): // abort if context signals done
		return false
	case m.box <- struct{}{}: // try to put something into the box
		return true
	}
}

func (m RWMutex) Unlock() {
	<-m.box // take something from the box
}

func (m RWMutex) RLock(ctx context.Context) bool {
	var currentReaders int32
	select {
	case <-ctx.Done():
		return false
	case m.box <- struct{}{}: // try to put something into the box
		// If the box is available we have no readers and no writers.
	case currentReaders = <-m.readers: // try to get the current number of readers
		// There are some readers. We need to update the count.
	}
	// Updated the readers count.
	currentReaders++
	m.readers <- currentReaders
	return true
}

func (m RWMutex) RUnlock() {
	// Take the value of readers and decrement it.
	currentReaders := <-m.readers
	currentReaders--
	if currentReaders == 0 {
		// We were the last reader. Make the box available again and return.
		<-m.box
		return
	}
	// Update the readers count, keep the box non-empty.
	m.readers <- currentReaders
}
