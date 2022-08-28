package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestBackoff_Do(t *testing.T) {
	type C struct {
		options             *Options
		expectedErrStr      string
		expectedErr         error
		expectedElapsedTime time.Duration
		expectedStat        string

		b *Backoff
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		t.Parallel()

		var options Options
		c.options = &options

		testcase.DoCallback(0, t, c)

		clock := clock.NewMock()
		clock.Set(time.Unix(123, 456))
		b := DoNew(clock, options)
		c.b = b

		testcase.DoCallback(1, t, c)

		t0 := time.Now()
		err := b.Do()
		elapsedTime := time.Since(t0)
		if c.expectedErrStr == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.expectedErrStr)
			if c.expectedErr != nil {
				assert.ErrorIs(t, err, c.expectedErr)
			}
		}
		if t.Failed() {
			return
		}
		assert.InDelta(t, c.expectedElapsedTime, elapsedTime, float64(10*time.Millisecond))
		assert.Equal(t, c.expectedStat, b.DumpAsString())
	}).
		SetCallback(1, func(*testing.T, *C) {})

	// MinDelay
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2.5)
		c.options.MaxDelayJitter.Set(0)
		c.expectedElapsedTime = 100 * time.Millisecond
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2.5
maxDelayJitter: 0
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 1
delay: 100ms
`[1:]
	}).Run(t)

	// DelayFactor
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2.5)
		c.options.MaxDelayJitter.Set(0)
		c.expectedElapsedTime = 250 * time.Millisecond
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2.5
maxDelayJitter: 0
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 2
delay: 250ms
`[1:]
	}).SetCallback(1, func(t *testing.T, c *C) {
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
	}).Run(t)

	// MaxDelay
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2.5)
		c.options.MaxDelayJitter.Set(0)
		c.expectedElapsedTime = 300 * time.Millisecond
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2.5
maxDelayJitter: 0
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 3
delay: 300ms
`[1:]
	}).SetCallback(1, func(t *testing.T, c *C) {
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
	}).Run(t)

	// MaxDelayJitter (1)
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2)
		c.options.MaxDelayJitter.Set(0.8)
		c.expectedElapsedTime = 157872371
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2
maxDelayJitter: 0.8
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 1
delay: 100ms
`[1:]
	}).Run(t)

	// MaxDelayJitter (2)
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2)
		c.options.MaxDelayJitter.Set(0.8)
		c.expectedElapsedTime = 183767889
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2
maxDelayJitter: 0.8
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 2
delay: 200ms
`[1:]
	}).SetCallback(1, func(t *testing.T, c *C) {
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
	}).Run(t)

	// MaxNumberOfAttempts
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(300 * time.Millisecond)
		c.options.DelayFactor.Set(2)
		c.options.MaxDelayJitter.Set(0)
		c.options.MaxNumberOfAttempts.Set(1)
		c.expectedErrStr = "backoff: too many attempts; maxNumberOfAttempts=1"
		c.expectedErr = ErrTooManyAttempts
		c.expectedElapsedTime = 0
		c.expectedStat = `
minDelay: 100ms
maxDelay: 300ms
delayFactor: 2
maxDelayJitter: 0
delayFuncIsNil: false
maxNumberOfAttempts: 1
attemptCount: 2
delay: 100ms
`[1:]
	}).SetCallback(1, func(t *testing.T, c *C) {
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
	}).Run(t)

	// DelayWithContext
	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.options.MinDelay.Set(100 * time.Millisecond)
		c.options.MaxDelay.Set(500 * time.Millisecond)
		c.options.DelayFactor.Set(2)
		c.options.MaxDelayJitter.Set(0)
		ctx, cancel := context.WithTimeout(context.Background(), 450*time.Millisecond)
		t.Cleanup(cancel)
		c.options.DelayFunc = DelayWithContext(ctx)
		c.expectedErrStr = "context deadline exceeded"
		c.expectedErr = context.DeadlineExceeded
		c.expectedElapsedTime = 150 * time.Millisecond
		c.expectedStat = `
minDelay: 100ms
maxDelay: 500ms
delayFactor: 2
maxDelayJitter: 0
delayFuncIsNil: false
maxNumberOfAttempts: 100
attemptCount: 3
delay: 400ms
`[1:]
	}).SetCallback(1, func(t *testing.T, c *C) {
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
		if err := c.b.Do(); err != nil {
			t.Fatal(err)
		}
	}).Run(t)
}
