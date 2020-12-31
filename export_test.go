package backoff

func (o *Options) Sanitize() { o.sanitize() }

func (b *Backoff) Options() *Options { return &b.options }
