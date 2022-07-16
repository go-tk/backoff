package backoff

import (
	"bytes"
	"fmt"
)

func (o *Options) Apply(to *Backoff) { o.apply(to) }

func (b *Backoff) DumpSettings() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "minDelay: %v\n", b.minDelay)
	fmt.Fprintf(&buffer, "maxDelay: %v\n", b.maxDelay)
	fmt.Fprintf(&buffer, "delayFactor: %v\n", b.delayFactor)
	fmt.Fprintf(&buffer, "maxDelayJitter: %v\n", b.maxDelayJitter)
	fmt.Fprintf(&buffer, "maxNumberOfAttempts: %v\n", b.maxNumberOfAttempts)
	return buffer.String()
}
