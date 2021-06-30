package backoff

import (
	"errors"
	"fmt"
)

// Backoff represents an instance of the exponential backoff algorithm.
type Backoff struct {
	options      Options
	attemptCount int
	delay        delay
	timer        timer
}

// New creates an instance of the exponential backoff algorithm with the
// given options.
func New(options Options) *Backoff {
	var b Backoff
	b.options = options
	return &b
}

// Do delays for a time period determined based on the options.
func (b *Backoff) Do() error {
	attemptCount := b.attemptCount
	b.attemptCount++
	if attemptCount == 0 {
		b.options.sanitize()
	}
	if b.options.MaxNumberOfAttempts.Value() >= 0 && attemptCount >= b.options.MaxNumberOfAttempts.Value() {
		return fmt.Errorf("%w; maxNumberOfAttempts=%v", ErrTooManyAttempts, b.options.MaxNumberOfAttempts)
	}
	delay := b.delay.Update(b.options.MinDelay.Value(), b.options.MaxDelay.Value(),
		b.options.DelayFactor.Value(), b.options.MaxDelayJitter.Value())
	timedOut := b.timer.Start(delay)
	defer b.timer.Stop()
	if err := b.options.DelayFunc(timedOut); err != nil {
		return err
	}
	return nil
}

// ErrTooManyAttempts is returned when the maximum number of attempts
// to back off has been reached.
var ErrTooManyAttempts = errors.New("backoff: too many attempts")
