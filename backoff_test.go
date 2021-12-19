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
	type Init struct {
		Options Options
	}
	type Output struct {
		Err error
	}
	type Workspace struct {
		B                *Backoff
		Init             Init
		ExpectedOutput   Output
		DurationInterval [2]time.Duration
	}
	tc := testcase.New().
		Step(1, func(t *testing.T, w *Workspace) {
			w.B = New(w.Init.Options)
		}).
		Step(2, func(t *testing.T, w *Workspace) {
			t0 := time.Now()
			var err error
			for {
				err = w.B.Do()
				if err != nil {
					break
				}
			}
			duration := time.Since(t0)
			for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
			}
			var output Output
			output.Err = err
			assert.Equal(t, w.ExpectedOutput, output)
			if w.DurationInterval != [2]time.Duration{} {
				assert.GreaterOrEqual(t, duration, w.DurationInterval[0])
				assert.LessOrEqual(t, duration, w.DurationInterval[1])
			}
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("MinDelay option").
			Then("should respect MinDelay option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(500 * time.Millisecond)
				w.Init.Options.MaxDelayJitter.Set(0)
				w.Init.Options.MaxNumberOfAttempts.Set(1)
				w.ExpectedOutput.Err = ErrTooManyAttempts
				w.DurationInterval[0] = 500 * time.Millisecond
				w.DurationInterval[1] = (500 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("MaxDelay option").
			Then("should respect MaxDelay option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(150 * time.Millisecond)
				w.Init.Options.MaxDelay.Set(200 * time.Millisecond)
				w.Init.Options.DelayFactor.Set(100)
				w.Init.Options.MaxDelayJitter.Set(0)
				w.Init.Options.MaxNumberOfAttempts.Set(3)
				w.ExpectedOutput.Err = ErrTooManyAttempts
				w.DurationInterval[0] = 550 * time.Millisecond
				w.DurationInterval[1] = (550 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("DelayFactor option").
			Then("should respect DelayFactor option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(50 * time.Millisecond)
				w.Init.Options.MaxDelay.Set(time.Hour)
				w.Init.Options.DelayFactor.Set(3)
				w.Init.Options.MaxDelayJitter.Set(0)
				w.Init.Options.MaxNumberOfAttempts.Set(3)
				w.ExpectedOutput.Err = ErrTooManyAttempts
				w.DurationInterval[0] = 650 * time.Millisecond
				w.DurationInterval[1] = (650 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("MaxDelayJitter option").
			Then("should respect MaxDelayJitter option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(400 * time.Millisecond)
				w.Init.Options.MaxDelayJitter.Set(0.25)
				w.Init.Options.MaxNumberOfAttempts.Set(1)
				w.ExpectedOutput.Err = ErrTooManyAttempts
				w.DurationInterval[0] = 300 * time.Millisecond
				w.DurationInterval[1] = (500 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("DelayFunc option").
			Then("should respect DelayFunc option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(time.Hour)
				w.Init.Options.MaxDelay.Set(time.Hour)
				w.Init.Options.MaxDelayJitter.Set(0)
				ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
				_ = cancel
				w.Init.Options.DelayFunc = DelayWithContext(ctx)
				w.Init.Options.MaxNumberOfAttempts.Set(1)
				w.ExpectedOutput.Err = context.DeadlineExceeded
				w.DurationInterval[0] = 300 * time.Millisecond
				w.DurationInterval[1] = (300 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("MaxNumberOfAttempts option with negative value").
			Then("should not limit number of attempts").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(50 * time.Millisecond)
				w.Init.Options.MaxDelay.Set(50 * time.Millisecond)
				w.Init.Options.MaxDelayJitter.Set(0)
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				_ = cancel
				w.Init.Options.DelayFunc = DelayWithContext(ctx)
				w.Init.Options.MaxNumberOfAttempts.Set(-1)
				w.ExpectedOutput.Err = context.DeadlineExceeded
				w.DurationInterval[0] = 500 * time.Millisecond
				w.DurationInterval[1] = (500 + 100) * time.Millisecond
			}),
	)
}
