package backoff_test

import (
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestOptions_Sanitize(t *testing.T) {
	const (
		DefaultMinDelay            time.Duration = 100 * time.Millisecond
		DefaultMaxDelay            time.Duration = 100 * time.Second
		DefaultDelayFactor         float64       = 2
		DefaultMaxDelayJitter      float64       = 1
		DefaultMaxNumberOfAttempts               = 100
	)
	type State struct {
		MinDelay            time.Duration
		MaxDelay            time.Duration
		DelayFactor         float64
		MaxDelayJitter      float64
		DelayFuncIsNil      bool
		MaxNumberOfAttempts int
	}
	type Context struct {
		O Options

		ExpectedState State
	}
	tc := testcase.New(func(t *testing.T) *Context {
		return &Context{
			ExpectedState: State{
				MinDelay:            DefaultMinDelay,
				MaxDelay:            DefaultMaxDelay,
				DelayFactor:         DefaultDelayFactor,
				MaxDelayJitter:      DefaultMaxDelayJitter,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			},
		}
	}).Run(func(t *testing.T, c *Context) {
		c.O.Sanitize()
		var state State
		state.MinDelay = c.O.MinDelay
		state.MaxDelay = c.O.MaxDelay
		state.DelayFactor = c.O.DelayFactor
		state.MaxDelayJitter = c.O.MaxDelayJitter
		state.DelayFuncIsNil = c.O.DelayFunc == nil
		state.MaxNumberOfAttempts = c.O.MaxNumberOfAttempts
		assert.Equal(t, c.ExpectedState, state)
	})
	testcase.RunListParallel(t,
		tc.Copy().
			Then("should overwrite options unset to default"),
		tc.Copy().
			Given("negative option MinDelay, MaxDelay, DelayFactor").
			Then("should overwrite options to default").
			PreRun(func(t *testing.T, c *Context) {
				c.O.MinDelay = -1
				c.O.MaxDelay = -2
				c.O.DelayFactor = -3
			}),
		tc.Copy().
			Given("options with specified values").
			Then("should not overwrite options").
			PreRun(func(t *testing.T, c *Context) {
				c.O.MinDelay = 1
				c.O.MaxDelay = 2
				c.O.DelayFactor = 3
				c.O.MaxDelayJitter = 4
				c.O.MaxNumberOfAttempts = 5
				c.ExpectedState = State{
					MinDelay:            1,
					MaxDelay:            2,
					DelayFactor:         3,
					MaxDelayJitter:      4,
					MaxNumberOfAttempts: 5,
				}
			}),
		tc.Copy().
			Given("option MaxDelay less than option MinDelay").
			Then("should overwrite option MaxDelay to MinDelay").
			PreRun(func(t *testing.T, c *Context) {
				c.O.MinDelay = 2
				c.O.MaxDelay = 1
				c.ExpectedState.MinDelay = 2
				c.ExpectedState.MaxDelay = 2
			}),
		tc.Copy().
			Given("negative option MaxDelayJitter, MaxNumberOfAttempts").
			Then("should overwrite option MaxDelayJitter to zero").
			PreRun(func(t *testing.T, c *Context) {
				c.O.MaxDelayJitter = -100
				c.O.MaxNumberOfAttempts = -100
				c.ExpectedState.MaxDelayJitter = 0
				c.ExpectedState.MaxNumberOfAttempts = 0
			}),
	)
}
