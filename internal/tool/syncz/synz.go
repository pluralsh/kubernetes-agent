package syncz

import "sync"

func RunWithMutex[T any](mu *sync.Mutex, f func() T) T {
	mu.Lock()
	defer mu.Unlock()
	return f()
}
