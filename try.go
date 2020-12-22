package try

import (
	"context"
	"math/rand"
	"time"
)

// Do tries the given function with the given options.
//
// When the given function returns true, it returns true as well.
//
// When the given function returns false,
//
// a) if option MaxNumberOfAttempts is set to positive and the number of attempts
// reaches the value of the option, it returns false as well.
//
// b) otherwise it waits for a backoff time, with respect to the backoff options,
// and then retry the function.
func Do(ctx context.Context, f func() (bool, error), options Options) (bool, error) {
	var backoffTimer timer
	for attemptCount := 1; ; attemptCount++ {
		ok, err := f()
		if err != nil {
			return ok, err
		}
		if ok {
			return true, nil
		}
		if attemptCount == 1 {
			options.normalize()
		}
		if attemptCount == options.MaxNumberOfAttempts {
			return false, nil
		}
		backoffTimer.Set(
			options.MinBackoff,
			options.MaxBackoff,
			options.BackoffFactor,
			options.MaxBackoffJitter,
		)
		select {
		case <-backoffTimer.C():
		case <-ctx.Done():
			backoffTimer.Stop()
			return false, ctx.Err()
		}
	}
}

type timer struct {
	t     *time.Timer
	rand  *rand.Rand
	delay time.Duration
}

func (t *timer) Set(minDelay, maxDelay time.Duration, delayFactor, maxDelayJitter float64) {
	if t.t == nil {
		t.delay = minDelay
	} else {
		t.delay = time.Duration(float64(t.delay) * delayFactor)
		if t.delay > maxDelay {
			t.delay = maxDelay
		}
	}
	if t.t == nil && maxDelayJitter > 0 {
		t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	delayWithJitter := t.delayWithJitter(maxDelayJitter)
	if t.t == nil {
		t.t = time.NewTimer(delayWithJitter)
	} else {
		t.t.Reset(delayWithJitter)
	}
}

func (t *timer) delayWithJitter(maxDelayJitter float64) time.Duration {
	if maxDelayJitter == 0 {
		return t.delay
	}
	delayJitter := maxDelayJitter * (2*t.rand.Float64() - 1)
	delayWithJitter := t.delay + time.Duration(float64(t.delay)*delayJitter)
	return delayWithJitter
}

func (t *timer) C() <-chan time.Time { return t.t.C }
func (t *timer) Stop()               { t.t.Stop() }
