package backoff

import (
	"bytes"
	"fmt"
)

var DoNew = doNew

func (o *Options) Apply(to *Backoff) { o.apply(to) }

func (b *Backoff) Dump(buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "minDelay: %v\n", b.minDelay)
	fmt.Fprintf(buffer, "maxDelay: %v\n", b.maxDelay)
	fmt.Fprintf(buffer, "delayFactor: %v\n", b.delayFactor)
	fmt.Fprintf(buffer, "maxDelayJitter: %v\n", b.maxDelayJitter)
	fmt.Fprintf(buffer, "delayFuncIsNil: %v\n", b.delayFunc == nil)
	fmt.Fprintf(buffer, "maxNumberOfAttempts: %v\n", b.maxNumberOfAttempts)
	fmt.Fprintf(buffer, "attemptCount: %v\n", b.attemptCount)
	fmt.Fprintf(buffer, "delay: %v\n", b.delay.d)
}

func (b *Backoff) DumpAsString() string {
	var buffer bytes.Buffer
	b.Dump(&buffer)
	return buffer.String()
}
