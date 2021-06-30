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

const (
	defaultMinDelay            = 100 * time.Millisecond
	defaultMaxDelay            = 100 * time.Second
	defaultDelayFactor         = 2
	defaultMaxDelayJitter      = 1
	defaultMaxNumberOfAttempts = 100
)

func (o *Options) sanitize() {
	if o.MinDelay.Value() <= 0 {
		o.MinDelay.Set(defaultMinDelay)
	}
	if o.MaxDelay.Value() <= 0 {
		o.MaxDelay.Set(defaultMaxDelay)
	}
	if o.MinDelay.Value() > o.MaxDelay.Value() {
		o.MinDelay.Set(defaultMinDelay)
		o.MaxDelay.Set(defaultMaxDelay)
	}
	if o.DelayFactor.Value() < 1 {
		o.DelayFactor.Set(defaultDelayFactor)
	}
	if !o.MaxDelayJitter.HasValue() {
		o.MaxDelayJitter.Set(defaultMaxDelayJitter)
	}
	if o.DelayFunc == nil {
		o.DelayFunc = func(timedOut <-chan struct{}) error {
			<-timedOut
			return nil
		}
	}
	if !o.MaxNumberOfAttempts.HasValue() {
		o.MaxNumberOfAttempts.Set(defaultMaxNumberOfAttempts)
	}
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
