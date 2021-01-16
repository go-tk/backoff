package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestBackoff_Do(t *testing.T) {
	const Delta = float64(40 * time.Millisecond)
	type Init struct {
		Options Options
	}
	type Output struct {
		Err error
	}
	type Context struct {
		B  *Backoff
		T0 time.Time

		Init           Init
		Output         Output
		ExceptedOutput Output
	}
	tc := testcase.New(func(t *testing.T) *Context {
		return &Context{}
	}).Setup(func(t *testing.T, c *Context) {
		c.B = New(c.Init.Options)
	}).Run(func(t *testing.T, c *Context) {
		var err error
		for {
			err = c.B.Do()
			if err != nil {
				break
			}
		}
		var output Output
		for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
		}
		output.Err = err
		assert.Equal(t, c.ExceptedOutput, output)
	})
	for i := 0; i < 10; i++ {
		testcase.RunList(t, []testcase.TestCase{
			tc.Copy().
				Given("option MinDelay").
				Then("should respect option MinDelay").
				PreSetup(func(t *testing.T, c *Context) {
					c.Init.Options.MinDelay = 200 * time.Millisecond
					c.Init.Options.MaxDelayJitter = -1
					c.Init.Options.MaxNumberOfAttempts = 1
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = ErrTooManyAttempts
					c.T0 = time.Now()
				}).
				PostRun(func(t *testing.T, c *Context) {
					d := time.Since(c.T0)
					assert.InDelta(t, 200*time.Millisecond, d, Delta)
				}),
			tc.Copy().
				Given("option MaxDelay").
				Then("should respect option MaxDelay").
				PreSetup(func(t *testing.T, c *Context) {
					c.Init.Options.MinDelay = 10 * time.Millisecond
					c.Init.Options.MaxDelay = 200 * time.Millisecond
					c.Init.Options.DelayFactor = 100
					c.Init.Options.MaxDelayJitter = -1
					c.Init.Options.MaxNumberOfAttempts = 2
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = ErrTooManyAttempts
					c.T0 = time.Now()
				}).
				PostRun(func(t *testing.T, c *Context) {
					d := time.Since(c.T0)
					assert.InDelta(t, 210*time.Millisecond, d, Delta)
				}),
			tc.Copy().
				Given("option DelayFactor").
				Then("should respect option DelayFactor").
				PreSetup(func(t *testing.T, c *Context) {
					c.Init.Options.MinDelay = 100 * time.Millisecond
					c.Init.Options.MaxDelay = 1 * time.Second
					c.Init.Options.DelayFactor = 1.5
					c.Init.Options.MaxDelayJitter = -1
					c.Init.Options.MaxNumberOfAttempts = 3
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = ErrTooManyAttempts
					c.T0 = time.Now()
				}).
				PostRun(func(t *testing.T, c *Context) {
					d := time.Since(c.T0)
					assert.InDelta(t, 475*time.Millisecond, d, Delta)
				}),
			tc.Copy().
				Given("option MaxDelayJitter").
				Then("should respect option MaxDelayJitter").
				PreSetup(func(t *testing.T, c *Context) {
					c.Init.Options.MinDelay = 200 * time.Millisecond
					c.Init.Options.MaxDelay = 200 * time.Second
					c.Init.Options.MaxDelayJitter = 0.3
					c.Init.Options.MaxNumberOfAttempts = 1
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = ErrTooManyAttempts
					c.T0 = time.Now()
				}).
				PostRun(func(t *testing.T, c *Context) {
					d := time.Since(c.T0)
					assert.InDelta(t, 200*time.Millisecond, d, float64(60*time.Millisecond)+Delta)
				}),
			tc.Copy().
				Given("option DelayFunc").
				Then("should respect option DelayFunc").
				PreSetup(func(t *testing.T, c *Context) {
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
					_ = cancel
					c.Init.Options.DelayFunc = DelayWithContext(ctx)
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = context.DeadlineExceeded
				}),
			tc.Copy().
				Given("option MaxNumberOfAttempts").
				Then("should respect option MaxNumberOfAttempts").
				PreSetup(func(t *testing.T, c *Context) {
					c.Init.Options.MaxNumberOfAttempts = -1
				}).
				PreRun(func(t *testing.T, c *Context) {
					c.ExceptedOutput.Err = ErrTooManyAttempts
					c.T0 = time.Now()
				}),
		})
	}
}
