package backoff

import (
	"errors"
	"fmt"
)

// Backoff represents an instance of the exponential backoff algorithm.
type Backoff struct {
	options      Options
	attemptCount int
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
	if b.attemptCount == 0 {
		b.options.sanitize()
	} else {
		if b.options.MaxNumberOfAttempts >= 1 && b.attemptCount == b.options.MaxNumberOfAttempts {
			return fmt.Errorf("%w; maxNumberOfAttempts=%v", ErrTooManyAttempts, b.options.MaxNumberOfAttempts)
		}
	}
	b.attemptCount++
	b.timer.Start(b.options.MinDelay, b.options.MaxDelay, b.options.DelayFactor, b.options.MaxDelayJitter)
	defer b.timer.Stop()
	event := b.timer.Expiration()
	if err := b.options.DelayFunc(event); err != nil {
		return err
	}
	return nil
}

// ErrTooManyAttempts is returned when the maximum number of attempts
// to back off has been reached.
var ErrTooManyAttempts = errors.New("backoff: too many attempts")
