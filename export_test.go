package backoff

const (
	DefaultMinDelay            = defaultMinDelay
	DefaultMaxDelay            = defaultMaxDelay
	DefaultDelayFactor         = defaultDelayFactor
	DefaultMaxDelayJitter      = defaultMaxDelayJitter
	DefaultMaxNumberOfAttempts = defaultMaxNumberOfAttempts
)

func (o *Options) Sanitize() { o.sanitize() }
