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
	MinBackoff          time.Duration
	MaxBackoff          time.Duration
	BackoffFactor       float64
	MaxBackoffJitter    float64
	MaxNumberOfAttempts int
}

func (o *Options) normalize() {
	if o.MinBackoff <= 0 {
		o.MinBackoff = DefaultMinBackoff
	}
	if o.MaxBackoff <= 0 {
		o.MaxBackoff = DefaultMaxBackoff
	}
	if o.MaxBackoff < o.MinBackoff {
		o.MaxBackoff = o.MinBackoff
	}
	if o.BackoffFactor <= 0 {
		o.BackoffFactor = DefaultBackoffFactor
	}
	if o.MaxBackoffJitter <= 0 {
		o.MaxBackoffJitter = DefaultMaxBackoffJitter
	}
	if o.MaxNumberOfAttempts == 0 {
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
