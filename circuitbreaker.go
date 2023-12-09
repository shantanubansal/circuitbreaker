package circuitbreaker

import (
	"sync"
	"time"
)

type SimpleCircuitBreaker struct {
	failureThreshold   int
	intervalDuration   time.Duration
	cooldownDuration   time.Duration
	operationCounters  *sync.Map
	lastOperationTimes *sync.Map
	mutex              *sync.Mutex
}

func NewSimpleCircuitBreaker(failureThreshold int, intervalDuration, cooldownDuration time.Duration) *SimpleCircuitBreaker {
	return &SimpleCircuitBreaker{
		failureThreshold:   failureThreshold,
		intervalDuration:   intervalDuration,
		cooldownDuration:   cooldownDuration,
		operationCounters:  new(sync.Map),
		lastOperationTimes: new(sync.Map),
		mutex:              &sync.Mutex{},
	}
}

func (cb *SimpleCircuitBreaker) AllowOperation(key string) bool {
	return cb.AllowOperationWithCustomCount(key, 1)
}

func (cb *SimpleCircuitBreaker) AllowOperationWithCustomCount(key string, executionCount int) bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	now := time.Now()
	if lastTime, ok := cb.lastOperationTimes.Load(key); ok {
		if now.Sub(lastTime.(time.Time)) > cb.cooldownDuration {
			cb.operationCounters.Delete(key)
			cb.lastOperationTimes.Delete(key)
		}
	}

	counter, _ := cb.operationCounters.LoadOrStore(key, 0)
	if counter.(int) < cb.failureThreshold {
		cb.operationCounters.Store(key, counter.(int)+executionCount)
		return true
	}

	if _, ok := cb.lastOperationTimes.Load(key); !ok {
		cb.lastOperationTimes.Store(key, now)
	}
	return false
}
