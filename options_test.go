package backoff_test

import (
	"context"
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestOptions_Apply(t *testing.T) {
	type C struct {
		options      *Options
		expectedStat string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		var options Options
		c.options = &options

		testcase.DoCallback(0, t, c)

		var b Backoff
		options.Apply(&b)
		assert.Equal(t, c.expectedStat, b.DumpAsString())
	})

	// default option values
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.expectedStat = `
minDelay: 100ms
maxDelay: 1m40s
delayFactor: 2
maxDelayJitter: 1
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 0
delay: 0s
`[1:]
	}).Run(t)

	// invalid option values (1)
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(-1 * time.Second)
		c.options.MaxDelay.Set(-1 * time.Second)
		c.options.DelayFactor.Set(0.5)
		c.options.MaxDelayJitter.Set(1.1)
		c.expectedStat = `
minDelay: 100ms
maxDelay: 1m40s
delayFactor: 2
maxDelayJitter: 1
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 0
delay: 0s
`[1:]
	}).Run(t)

	// invalid option values (2)
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(2 * time.Second)
		c.options.MaxDelay.Set(1 * time.Second)
		c.expectedStat = `
minDelay: 100ms
maxDelay: 1m40s
delayFactor: 2
maxDelayJitter: 1
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 0
delay: 0s
`[1:]
	}).Run(t)

	// valid option values
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(1 * time.Second)
		c.options.MaxDelay.Set(2 * time.Second)
		c.options.DelayFactor.Set(1.5)
		c.options.MaxDelayJitter.Set(0.8)
		c.options.DelayFunc = DelayWithContext(context.Background())
		c.options.MaxNumberOfAttempts.Set(10)
		c.expectedStat = `
minDelay: 1s
maxDelay: 2s
delayFactor: 1.5
maxDelayJitter: 0.8
delayFuncIsNil: false
maxNumberOfAttempts: 10
attemptCount: 0
delay: 0s
`[1:]
	}).Run(t)
}
