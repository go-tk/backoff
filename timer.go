package backoff

import (
	"math/rand"
	"time"
)

type timer struct {
	delay         time.Duration
	rand          *rand.Rand
	t             *time.Timer
	expiration    <-chan struct{}
	isInitialized bool
}

func (t *timer) Start(minDelay, maxDelay time.Duration, delayFactor, maxDelayJitter float64) {
	if t.isInitialized {
		t.delay = time.Duration(float64(t.delay) * delayFactor)
		if t.delay > maxDelay {
			t.delay = maxDelay
		}
	} else {
		t.delay = minDelay
		if maxDelayJitter > 0 {
			t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
	}
	delayWithJitter := t.delayWithJitter(maxDelayJitter)
	expiration := make(chan struct{})
	t.t = time.AfterFunc(delayWithJitter, func() { close(expiration) })
	t.expiration = expiration
	t.isInitialized = true
}

func (t *timer) delayWithJitter(maxDelayJitter float64) time.Duration {
	if maxDelayJitter == 0 {
		return t.delay
	}
	delayJitter := maxDelayJitter * (2*t.rand.Float64() - 1)
	delayWithJitter := t.delay + time.Duration(float64(t.delay)*delayJitter)
	return delayWithJitter
}

func (t *timer) Stop()                       { t.t.Stop() }
func (t *timer) Expiration() <-chan struct{} { return t.expiration }
