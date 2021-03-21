package backoff

import (
	"context"
	"time"
)

// Options represents options for backoffs.
//
// The pseudo code for the delay calculation:
//  if delay is not initialized
//      delay = MIN_DELAY
//  else
//      delay = min(delay * DELAY_FACTOR, MAX_DELAY)
//
//  delay_jitter = random(-MAX_DELAY_JITTER, MAX_DELAY_JITTER)
//  delay_with_jitter = delay * (1 + delay_jitter)
type Options struct {
	// Value < 1 is equivalent to 100ms.
	MinDelay time.Duration

	// Value < 1 is equivalent to 100s.
	MaxDelay time.Duration

	// Value < 1 is equivalent to 2.
	DelayFactor float64

	// Value < 0 is equivalent to 0.
	// Value == 0 is equivalent to 1.
	MaxDelayJitter float64

	// Value nil is equivalent to:
	//  func(event <-chan struct{}) error {
	//      <-event
	//      return nil
	//  }
	DelayFunc DelayFunc

	// Value < 0 is equivalent to no limit.
	// value == 0 is equivalent to 100.
	MaxNumberOfAttempts int
}

func (o *Options) sanitize() {
	if o.MinDelay < 1 {
		o.MinDelay = 100 * time.Millisecond
	}
	if o.MaxDelay < 1 {
		o.MaxDelay = 100 * time.Second
	}
	if o.MaxDelay < o.MinDelay {
		o.MaxDelay = o.MinDelay
	}
	if o.DelayFactor < 1 {
		o.DelayFactor = 2
	}
	if o.MaxDelayJitter < 0 {
		o.MaxDelayJitter = 0
	} else if o.MaxDelayJitter == 0 {
		o.MaxDelayJitter = 1
	}
	if o.DelayFunc == nil {
		o.DelayFunc = func(event <-chan struct{}) error {
			<-event
			return nil
		}
	}
	if o.MaxNumberOfAttempts < 0 {
		o.MaxNumberOfAttempts = 0
	} else if o.MaxNumberOfAttempts == 0 {
		o.MaxNumberOfAttempts = 100
	}
}

// DelayFunc is the type of the function delaying until the given event happens.
type DelayFunc func(event <-chan struct{}) (err error)

// DelayWithContext makes a DelayFunc with respect to the given ctx.
func DelayWithContext(ctx context.Context) DelayFunc {
	return func(event <-chan struct{}) error {
		select {
		case <-event:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
