package try_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	. "github.com/go-tk/try"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	type Input struct {
		Ctx     context.Context
		F       func() (bool, error)
		Options Options
	}
	type Output struct {
		OK     bool
		ErrStr string
	}
	type TestCase struct {
		Given, When, Then string
		Setup, Teardown   func(*TestCase)
		Input             Input
		Output            Output

		t0 time.Time
	}
	testCases := []TestCase{
		{
			Then: "should respect option MinBackoff",
			Setup: func(tc *TestCase) {
				tc.t0 = time.Now()
			},
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, nil },
				Options: Options{
					MinBackoff:          10 * time.Millisecond,
					MaxBackoffJitter:    -1,
					MaxNumberOfAttempts: 2,
				},
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.True(t, d >= 10*time.Millisecond)
				assert.True(t, d <= 20*time.Millisecond)
			},
		},
		{
			Then: "should respect option MaxBackoff",
			Setup: func(tc *TestCase) {
				tc.t0 = time.Now()
			},
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, nil },
				Options: Options{
					MinBackoff:          10 * time.Millisecond,
					MaxBackoff:          20 * time.Millisecond,
					BackoffFactor:       10,
					MaxBackoffJitter:    -1,
					MaxNumberOfAttempts: 3,
				},
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.True(t, d >= 30*time.Millisecond)
				assert.True(t, d <= 40*time.Millisecond)
			},
		},
		{
			Then: "should respect option BackoffFactor",
			Setup: func(tc *TestCase) {
				tc.t0 = time.Now()
			},
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, nil },
				Options: Options{
					MinBackoff:          10 * time.Millisecond,
					MaxBackoff:          30 * time.Millisecond,
					BackoffFactor:       3,
					MaxBackoffJitter:    -1,
					MaxNumberOfAttempts: 3,
				},
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.True(t, d >= 40*time.Millisecond)
				assert.True(t, d <= 50*time.Millisecond)
			},
		},
		{
			Then: "should respect option MaxBackoffJitter",
			Setup: func(tc *TestCase) {
				tc.t0 = time.Now()
			},
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, nil },
				Options: Options{
					MinBackoff:          10 * time.Millisecond,
					MaxBackoff:          20 * time.Millisecond,
					BackoffFactor:       2,
					MaxBackoffJitter:    1,
					MaxNumberOfAttempts: 3,
				},
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.True(t, d <= 70*time.Millisecond)
			},
		},
		{
			Then: "should respect option MaxNumberOfAttempts",
			Setup: func(tc *TestCase) {
				tc.t0 = time.Now()
			},
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, nil },
				Options: Options{
					MinBackoff:          10 * time.Millisecond,
					MaxBackoff:          20 * time.Millisecond,
					BackoffFactor:       2,
					MaxBackoffJitter:    -1,
					MaxNumberOfAttempts: 1,
				},
			},
			Teardown: func(tc *TestCase) {
				d := time.Since(tc.t0)
				assert.True(t, d <= 10*time.Millisecond)
			},
		},
		{
			When: "f returns true",
			Then: "should return true as well",
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return true, nil },
			},
			Output: Output{
				OK: true,
			},
		},
		{
			When: "f returns error",
			Then: "should return error as well",
			Input: Input{
				Ctx: context.Background(),
				F:   func() (bool, error) { return false, errors.New("my error") },
			},
			Output: Output{
				ErrStr: "my error",
			},
		},
		{
			When: "context is timed out",
			Then: "should return error",
			Setup: func(tc *TestCase) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
				_ = cancel
				tc.Input.Ctx = ctx
			},
			Input: Input{
				F: func() (bool, error) { return false, nil },
			},
			Output: Output{
				ErrStr: context.DeadlineExceeded.Error(),
			},
		},
	}
	for i := range testCases {
		tc := &testCases[i]
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			t.Logf("\nGIVEN: %s\nWHEN: %s\nTHEN: %s", tc.Given, tc.When, tc.Then)

			if f := tc.Setup; f != nil {
				f(tc)
			}

			ok, err := Do(tc.Input.Ctx, tc.Input.F, tc.Input.Options)

			var output Output
			output.OK = ok
			if err != nil {
				output.ErrStr = err.Error()
			}
			assert.Equal(t, tc.Output, output)

			if f := tc.Teardown; f != nil {
				f(tc)
			}
		})
	}
}
