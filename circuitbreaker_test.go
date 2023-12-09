package circuitbreaker_test

import (
	"github.com/shantanubansal/circuitbreaker"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	limiter := circuitbreaker.NewSimpleCircuitBreaker(5, 20*time.Second, 3*time.Second)
	key := "circuitbreaker-test"
	for i := 0; i < 8; i++ {
		isOpAllowed := limiter.AllowOperation(key)
		if i > 5 && isOpAllowed {
			t.Errorf("TestCircuitBreaker Failed %v, this iteration should have failed", i)
			return
		}
	}
	t.Logf("Waiting for 5 Seconds, so that circuit becomes active again")
	time.Sleep(5 * time.Second)
	if !limiter.AllowOperation(key) {
		t.Error("TestCircuitBreaker Failed: this iteration should have succeeded ")
	}
}
