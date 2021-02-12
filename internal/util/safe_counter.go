package util

import "sync"

type SafeCounter struct {
	lock  *sync.RWMutex
	count int
}

func NewSafeCounter() *SafeCounter {
	c := new(SafeCounter)
	c.init()
	return c
}

func (counter *SafeCounter) init() {
	counter.lock = new(sync.RWMutex)
	counter.count = 0
}

func (counter *SafeCounter) Increment() {
	counter.lock.Lock()
	counter.count++
	counter.lock.Unlock()
}

func (counter *SafeCounter) Decrement() {
	counter.lock.Lock()
	counter.count--
	counter.lock.Unlock()
}

func (counter *SafeCounter) Get() int {
	counter.lock.RLock()
	defer counter.lock.RUnlock()
	return counter.count
}
