package backoff_test

import (
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestOptions_Sanitize(t *testing.T) {
	type State struct {
		MinDelay            time.Duration
		MaxDelay            time.Duration
		DelayFactor         float64
		MaxDelayJitter      float64
		DelayFuncIsNil      bool
		MaxNumberOfAttempts int
	}
	type Workspace struct {
		testcase.WorkspaceBase

		O             Options
		ExpectedState State
	}
	tc := testcase.New().
		AddTask(1000, func(w *Workspace) {
			w.ExpectedState = State{
				MinDelay:            DefaultMinDelay,
				MaxDelay:            DefaultMaxDelay,
				DelayFactor:         DefaultDelayFactor,
				MaxDelayJitter:      DefaultMaxDelayJitter,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			}
		}).
		AddTask(2000, func(w *Workspace) {
			w.O.Sanitize()
			state := State{
				MinDelay:            w.O.MinDelay.Value(),
				MaxDelay:            w.O.MaxDelay.Value(),
				DelayFactor:         w.O.DelayFactor.Value(),
				MaxDelayJitter:      w.O.MaxDelayJitter.Value(),
				DelayFuncIsNil:      w.O.DelayFunc == nil,
				MaxNumberOfAttempts: w.O.MaxNumberOfAttempts.Value(),
			}
			assert.Equal(w.T(), w.ExpectedState, state)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("no option values").
			Then("should set option values to default"),
		tc.Copy().
			Given("invalid option values").
			Then("should set invalid option values to default").
			AddTask(1999, func(w *Workspace) {
				w.O.MinDelay.Set(-1)
				w.O.MaxDelay.Set(-1)
				w.O.DelayFactor.Set(-1)
			}),
		tc.Copy().
			Given("valid option values").
			Then("should preserve option values").
			AddTask(1999, func(w *Workspace) {
				w.O.MinDelay.Set(1 * time.Second)
				w.ExpectedState.MinDelay = w.O.MinDelay.Value()
				w.O.MaxDelay.Set(2 * time.Second)
				w.ExpectedState.MaxDelay = w.O.MaxDelay.Value()
				w.O.DelayFactor.Set(3)
				w.ExpectedState.DelayFactor = w.O.DelayFactor.Value()
				w.O.MaxDelayJitter.Set(0.3)
				w.ExpectedState.MaxDelayJitter = w.O.MaxDelayJitter.Value()
				w.O.MaxNumberOfAttempts.Set(0)
				w.ExpectedState.MaxNumberOfAttempts = w.O.MaxNumberOfAttempts.Value()
			}),
		tc.Copy().
			Given("MinDelay option value > MaxDelay option value").
			Then("should set MinDelay/MaxDelay option values to default").
			AddTask(1999, func(w *Workspace) {
				w.O.MinDelay.Set(2 * time.Second)
				w.O.MaxDelay.Set(1 * time.Second)
			}),
	)
}
