package backoff

import (
	"errors"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
)

// Backoff represents an instance of the exponential backoff algorithm.
type Backoff struct {
	clock clock.Clock

	minDelay            time.Duration
	maxDelay            time.Duration
	delayFactor         float64
	maxDelayJitter      float64
	delayFunc           DelayFunc
	maxNumberOfAttempts int

	attemptCount int
	delay        delay
	timer        timer
}

// New creates an instance of the exponential backoff algorithm with the
// given options.
func New(options Options) *Backoff { return doNew(clock.New(), options) }

func doNew(clock clock.Clock, options Options) *Backoff {
	var b Backoff
	b.clock = clock
	options.apply(&b)
	return &b
}

// Do delays for a time period determined based on the options.
func (b *Backoff) Do() error {
	attemptCount := b.attemptCount
	b.attemptCount++
	if b.maxNumberOfAttempts >= 0 && attemptCount >= b.maxNumberOfAttempts {
		return fmt.Errorf("%w; maxNumberOfAttempts=%v", ErrTooManyAttempts, b.maxNumberOfAttempts)
	}
	delay := b.delay.Update(b.minDelay, b.maxDelay, b.delayFactor, b.maxDelayJitter, b.clock)
	timedOut := b.timer.Start(delay)
	defer b.timer.Stop()
	if err := b.delayFunc(timedOut); err != nil {
		return err
	}
	return nil
}

// ErrTooManyAttempts is returned when the maximum number of attempts
// to back off has been reached.
var ErrTooManyAttempts = errors.New("backoff: too many attempts")
