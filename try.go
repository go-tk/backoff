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
func Do(ctx context.Context, f func() bool, options Options) (bool, error) {
	var backoff time.Duration
	var rand1 *rand.Rand
	var timer *time.Timer
	for attemptCount := 1; ; attemptCount++ {
		if f() {
			return true, nil
		}
		if attemptCount == 1 {
			options.normalize()
		}
		if attemptCount == options.MaxNumberOfAttempts {
			return false, nil
		}
		updateBackoff(&backoff, &options, &rand1, &timer)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return false, ctx.Err()
		}
	}
}

func updateBackoff(backoff *time.Duration, options *Options, rand1 **rand.Rand, timer **time.Timer) {
	if *timer == nil {
		*backoff = options.MinBackoff
	} else {
		*backoff = time.Duration(float64(*backoff) * options.BackoffFactor)
		if *backoff > options.MaxBackoff {
			*backoff = options.MaxBackoff
		}
	}
	if *timer == nil && options.MaxBackoffJitter > 0 {
		*rand1 = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	backoffWithJitter := makeBackoffWithJitter(*backoff, options.MaxBackoffJitter, *rand1)
	if *timer == nil {
		*timer = time.NewTimer(backoffWithJitter)
	} else {
		(*timer).Reset(backoffWithJitter)
	}
}

func makeBackoffWithJitter(backoff time.Duration, maxBackoffJitter float64, rand1 *rand.Rand) time.Duration {
	if maxBackoffJitter == 0 {
		return backoff
	}
	backoffJitter := maxBackoffJitter * (2*rand1.Float64() - 1)
	backoffWithJitter := backoff + time.Duration(float64(backoff)*backoffJitter)
	return backoffWithJitter
}
