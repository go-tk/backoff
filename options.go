package backoff

import (
	"context"
	"time"

	"github.com/go-tk/optional"
)

// Options represents options for backoffs.
//
// The pseudo code for the delay calculation:
//  if delay is not initialized
//      delay = MIN_DELAY
//  else
//      delay = min(delay * DELAY_FACTOR, MAX_DELAY)
//  delay_jitter = random(-MAX_DELAY_JITTER, MAX_DELAY_JITTER)
//  delay_with_jitter = delay * (1 + delay_jitter)
type Options struct {
	// Default value is 100ms.
	// Value <= 0 is equivalent to default value.
	MinDelay optional.Duration

	// Default value is 100s.
	// Value <= 0 is equivalent to default value.
	MaxDelay optional.Duration

	// Default value is 2.
	// Value < 1 is equivalent to default value.
	DelayFactor optional.Float64

	// Default value is 1.
	// Value <= 0 means no delay jitter.
	MaxDelayJitter optional.Float64

	// Default value is:
	//  func(timedOut <-chan struct{}) error {
	//      <-timedOut
	//      return nil
	//  }
	DelayFunc DelayFunc

	// Default value is 100.
	// Value < 0 means no limit on the number of attempts.
	MaxNumberOfAttempts optional.Int
}

func (o *Options) apply(to *Backoff) {
	if v, ok := o.MinDelay.Get(); ok && v > 0 {
		to.minDelay = v
	} else {
		to.minDelay = defaultMinDelay
	}
	if v, ok := o.MaxDelay.Get(); ok && v > 0 {
		to.maxDelay = v
	} else {
		to.maxDelay = defaultMaxDelay
	}
	if to.minDelay > to.maxDelay {
		to.minDelay = defaultMinDelay
		to.maxDelay = defaultMaxDelay
	}
	if v, ok := o.DelayFactor.Get(); ok && v >= 1 {
		to.delayFactor = v
	} else {
		to.delayFactor = defaultDelayFactor
	}
	if v, ok := o.MaxDelayJitter.Get(); ok {
		to.maxDelayJitter = v
	} else {
		to.maxDelayJitter = defaultMaxDelayJitter
	}
	if o.DelayFunc == nil {
		to.delayFunc = defaultDelayFunc
	} else {
		to.delayFunc = o.DelayFunc
	}
	if v, ok := o.MaxNumberOfAttempts.Get(); ok {
		to.maxNumberOfAttempts = v
	} else {
		to.maxNumberOfAttempts = defaultMaxNumberOfAttempts
	}
}

const (
	defaultMinDelay            = 100 * time.Millisecond
	defaultMaxDelay            = 100 * time.Second
	defaultDelayFactor         = 2
	defaultMaxDelayJitter      = 1
	defaultMaxNumberOfAttempts = 100
)

func defaultDelayFunc(timedOut <-chan struct{}) error {
	<-timedOut
	return nil
}

// DelayFunc is the type of the function blocking until the timed-out event happens.
type DelayFunc func(timedOut <-chan struct{}) (err error)

// DelayWithContext makes a DelayFunc with respect to the given ctx.
func DelayWithContext(ctx context.Context) DelayFunc {
	return func(timedOut <-chan struct{}) error {
		select {
		case <-timedOut:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
