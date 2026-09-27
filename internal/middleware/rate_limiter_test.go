package middleware

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
rl := NewRateLimiter(3 ,time.Minute)
key := "user-1"

if !rl.Allow(key){
	t.Error("request 1 should be allowed")
}
	if !rl.Allow(key) {
		t.Error("request 2 should be allowed")
	}

	if !rl.Allow(key) {
		t.Error("request 3 should be allowed")
	}

	if rl.Allow(key) {
		t.Error("request 4 should be blocked")
	}
}

func TestRateLimiterWindowReset(t *testing.T) {
	rl := NewRateLimiter(1, 100*time.Millisecond)

	key := "user-2"

	if !rl.Allow(key) {
		t.Error("first request should be allowed")
	}

	if rl.Allow(key) {
		t.Error("second request should be blocked")
	}

	time.Sleep(150 * time.Millisecond)

	if !rl.Allow(key) {
		t.Error("request should be allowed after window reset")
	}
}

func TestRateLimiterConcurrent( t *testing.T){
		rl := NewRateLimiter(10, time.Minute)

		key := "user-3"

		var wg sync.WaitGroup

	allowed := 0
var mu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if rl.Allow(key) {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if allowed != 10 {
		t.Errorf("expected 10 allowed requests, got %d", allowed)
	}

}