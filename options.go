package try

import "time"

// Options represents options for tries.
//
// The pseudo code for the backoff calculation:
//  if backoff is not initialized
//      backoff = MIN_BACKOFF
//  else
//      backoff = min(backoff * BACKOFF_FACTOR, MAX_BACKOFF)
//
//  backoff_jitter = random(-MAX_BACKOFF_JITTER, MAX_BACKOFF_JITTER)
//  backoff_with_jitter = backoff * (1 + backoff_jitter)
type Options struct {
	// Value < 1 for MinBackoff is equivalent to DefaultMinBackoff.
	MinBackoff time.Duration

	// Value < 1 for MaxBackoff is equivalent to DefaultMaxBackoff.
	MaxBackoff time.Duration

	// Value < 1 for BackoffFactor is equivalent to DefaultBackoffFactor.
	BackoffFactor float64

	// Value < 1 for MaxBackoffJitter indicates no backoff jitter.
	// Value == 0 for MaxBackoffJitter is equivalent to DefaultMaxBackoffJitter.
	MaxBackoffJitter float64

	// Value < 1 for MaxNumberOfAttempts indicates the number of attempts is unlimited.
	// value == 0 for MaxNumberOfAttempts is equivalent to DefaultMaxNumberOfAttempts.
	MaxNumberOfAttempts int
}

func (o *Options) normalize() {
	if o.MinBackoff < 1 {
		o.MinBackoff = DefaultMinBackoff
	}
	if o.MaxBackoff < 1 {
		o.MaxBackoff = DefaultMaxBackoff
	}
	if o.MaxBackoff < o.MinBackoff {
		o.MaxBackoff = o.MinBackoff
	}
	if o.BackoffFactor < 1 {
		o.BackoffFactor = DefaultBackoffFactor
	}
	if o.MaxBackoffJitter < 0 {
		o.MaxBackoffJitter = 0
	} else if o.MaxBackoffJitter == 0 {
		o.MaxBackoffJitter = DefaultMaxBackoffJitter
	}
	if o.MaxNumberOfAttempts < 0 {
		o.MaxNumberOfAttempts = 0
	} else if o.MaxNumberOfAttempts == 0 {
		o.MaxNumberOfAttempts = DefaultMaxNumberOfAttempts
	}
}

// Default values for options.
var (
	DefaultMinBackoff          time.Duration = 100 * time.Millisecond
	DefaultMaxBackoff          time.Duration = 100 * time.Second
	DefaultBackoffFactor       float64       = 2
	DefaultMaxBackoffJitter    float64       = 1
	DefaultMaxNumberOfAttempts int
)
