package backoff_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/stretchr/testify/assert"
)

func TestBackoff_Do(t *testing.T) {
	type Output struct {
		Err error
	}
	type TestCase struct {
		Given, When, Then string
		Setup, Teardown   func(*TestCase)
		Output            Output

		t  *testing.T
		b  *Backoff
		t0 time.Time
	}
	const delta = float64(40 * time.Millisecond)
	testCases := []TestCase{
		{
			Given: "option MinDelay",
			Then:  "should respect option MinDelay",
			Setup: func(tc *TestCase) {
				o := tc.b.Options()
				o.MinDelay = 200 * time.Millisecond
				o.MaxDelayJitter = -1
				o.MaxNumberOfAttempts = 1
				tc.t0 = time.Now()
			},
			Output: Output{
				Err: ErrTooManyAttempts,
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.InDelta(tc.t, 200*time.Millisecond, d, delta)
			},
		},
		{
			Given: "option MaxDelay",
			Then:  "should respect option MaxDelay",
			Setup: func(tc *TestCase) {
				o := tc.b.Options()
				o.MinDelay = 10 * time.Millisecond
				o.MaxDelay = 200 * time.Millisecond
				o.DelayFactor = 100
				o.MaxDelayJitter = -1
				o.MaxNumberOfAttempts = 2
				tc.t0 = time.Now()
			},
			Output: Output{
				Err: ErrTooManyAttempts,
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.InDelta(tc.t, 210*time.Millisecond, d, delta)
			},
		},
		{
			Given: "option DelayFactor",
			Then:  "should respect option DelayFactor",
			Setup: func(tc *TestCase) {
				o := tc.b.Options()
				o.MinDelay = 100 * time.Millisecond
				o.MaxDelay = 1 * time.Second
				o.DelayFactor = 1.5
				o.MaxDelayJitter = -1
				o.MaxNumberOfAttempts = 3
				tc.t0 = time.Now()
			},
			Output: Output{
				Err: ErrTooManyAttempts,
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.InDelta(tc.t, 475*time.Millisecond, d, delta)
			},
		},
		{
			Given: "option MaxDelayJitter",
			Then:  "should respect option MaxDelayJitter",
			Setup: func(tc *TestCase) {
				o := tc.b.Options()
				o.MinDelay = 200 * time.Millisecond
				o.MaxDelay = 200 * time.Second
				o.MaxDelayJitter = 0.3
				o.MaxNumberOfAttempts = 1
				tc.t0 = time.Now()
			},
			Output: Output{
				Err: ErrTooManyAttempts,
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.InDelta(tc.t, 200*time.Millisecond, d, float64(60*time.Millisecond)+delta)
			},
		},
		{
			Given: "option DelayFunc",
			Then:  "should respect option DelayFunc",
			Setup: func(tc *TestCase) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				_ = cancel
				o := tc.b.Options()
				o.DelayFunc = DelayWithContext(ctx)
			},
			Output: Output{
				Err: context.DeadlineExceeded,
			},
		},
		{
			Given: "option MaxNumberOfAttempts",
			Then:  "should respect option MaxNumberOfAttempts",
			Setup: func(tc *TestCase) {
				o := tc.b.Options()
				o.MaxNumberOfAttempts = -1
				tc.t0 = time.Now()
			},
			Output: Output{
				Err: ErrTooManyAttempts,
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.InDelta(tc.t, 0*time.Millisecond, d, delta)
				err := tc.b.Do()
				assert.True(t, errors.Is(err, ErrTooManyAttempts))
			},
		},
	}
	for i := 0; i < len(testCases)*10; i++ {
		tc := testCases[i%len(testCases)]
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			t.Logf("\nGIVEN: %s\nWHEN: %s\nTHEN: %s", tc.Given, tc.When, tc.Then)
			tc.t = t

			b := New(Options{})
			tc.b = b

			if f := tc.Setup; f != nil {
				f(&tc)
			}

			var err error
			for {
				err = b.Do()
				if err != nil {
					break
				}
			}

			var output Output
			for err2 := errors.Unwrap(err); err2 != nil; err, err2 = err2, errors.Unwrap(err2) {
			}
			output.Err = err
			assert.Equal(t, tc.Output, output)

			if f := tc.Teardown; f != nil {
				f(&tc)
			}
		})
	}
}
