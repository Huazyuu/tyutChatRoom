package utils

import "sync"

type ThreadSafety struct {
	mu sync.Mutex
}

var Safety ThreadSafety

func (receiver *ThreadSafety) Do(x func() interface{}) interface{} {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	return x()
}
