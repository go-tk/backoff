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
	type Workspace struct {
		B    *Backoff
		Init struct {
			Options Options
		}
		ExpOut, ActOut struct {
			Err error
		}
		Dur      time.Duration
		DurLimit [2]time.Duration
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
			w.Dur = time.Since(t0)
			for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
			}
			w.ActOut.Err = err
		}).
		Step(3, func(t *testing.T, w *Workspace) {
			assert.Equal(t, w.ExpOut, w.ActOut)
			if w.DurLimit != [2]time.Duration{} {
				assert.GreaterOrEqual(t, w.Dur, w.DurLimit[0])
				assert.LessOrEqual(t, w.Dur, w.DurLimit[1])
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
				w.ExpOut.Err = ErrTooManyAttempts
				w.DurLimit[0] = 500 * time.Millisecond
				w.DurLimit[1] = (500 + 100) * time.Millisecond
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
				w.ExpOut.Err = ErrTooManyAttempts
				w.DurLimit[0] = 550 * time.Millisecond
				w.DurLimit[1] = (550 + 100) * time.Millisecond
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
				w.ExpOut.Err = ErrTooManyAttempts
				w.DurLimit[0] = 650 * time.Millisecond
				w.DurLimit[1] = (650 + 100) * time.Millisecond
			}),
		tc.Copy().
			Given("MaxDelayJitter option").
			Then("should respect MaxDelayJitter option").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.Init.Options.MinDelay.Set(400 * time.Millisecond)
				w.Init.Options.MaxDelayJitter.Set(0.25)
				w.Init.Options.MaxNumberOfAttempts.Set(1)
				w.ExpOut.Err = ErrTooManyAttempts
				w.DurLimit[0] = 300 * time.Millisecond
				w.DurLimit[1] = (500 + 100) * time.Millisecond
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
				w.ExpOut.Err = context.DeadlineExceeded
				w.DurLimit[0] = 300 * time.Millisecond
				w.DurLimit[1] = (300 + 100) * time.Millisecond
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
				w.ExpOut.Err = context.DeadlineExceeded
				w.DurLimit[0] = 500 * time.Millisecond
				w.DurLimit[1] = (500 + 100) * time.Millisecond
			}),
	)
}
