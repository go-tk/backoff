package try_test

import (
	"strconv"
	"testing"

	. "github.com/go-tk/try"
	"github.com/stretchr/testify/assert"
)

func TestOptions_Sanitize(t *testing.T) {
	type State = Options
	type TestCase struct {
		Given, When, Then string
		Setup, Teardown   func(*TestCase)
		State             State

		o *Options
	}
	testCases := []TestCase{
		{
			Then: "should overwrite options to default",
			State: State{
				MinBackoff:          DefaultMinBackoff,
				MaxBackoff:          DefaultMaxBackoff,
				BackoffFactor:       DefaultBackoffFactor,
				MaxBackoffJitter:    DefaultMaxBackoffJitter,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MinBackoff, MaxBackoff, BackoffFactor",
			Then:  "should overwrite options to default",
			Setup: func(tc *TestCase) {
				tc.o.MinBackoff = -1
				tc.o.MaxBackoff = -2
				tc.o.BackoffFactor = -3
			},
			State: State{
				MinBackoff:          DefaultMinBackoff,
				MaxBackoff:          DefaultMaxBackoff,
				BackoffFactor:       DefaultBackoffFactor,
				MaxBackoffJitter:    DefaultMaxBackoffJitter,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "options with specified values",
			Then:  "should not overwrite options",
			Setup: func(tc *TestCase) {
				*tc.o = tc.State
			},
			State: State{
				MinBackoff:          1,
				MaxBackoff:          2,
				BackoffFactor:       3,
				MaxBackoffJitter:    4,
				MaxNumberOfAttempts: 5,
			},
		},
		{
			Given: "option MaxBackoff less than option MinBackoff",
			Then:  "should overwrite option MaxBackoff to MinBackoff",
			Setup: func(tc *TestCase) {
				tc.o.MinBackoff = 2
				tc.o.MaxBackoff = 1
			},
			State: State{
				MinBackoff:          2,
				MaxBackoff:          2,
				BackoffFactor:       DefaultBackoffFactor,
				MaxBackoffJitter:    DefaultMaxBackoffJitter,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MaxBackoffJitter",
			Then:  "should overwrite option MaxBackoffJitter to zero",
			Setup: func(tc *TestCase) {
				tc.o.MaxBackoffJitter = -100
			},
			State: State{
				MinBackoff:          DefaultMinBackoff,
				MaxBackoff:          DefaultMaxBackoff,
				BackoffFactor:       DefaultBackoffFactor,
				MaxNumberOfAttempts: DefaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MaxNumberOfAttempts",
			Then:  "should overwrite option MaxNumberOfAttempts to zero",
			Setup: func(tc *TestCase) {
				tc.o.MaxNumberOfAttempts = -100
			},
			State: State{
				MinBackoff:       DefaultMinBackoff,
				MaxBackoff:       DefaultMaxBackoff,
				BackoffFactor:    DefaultBackoffFactor,
				MaxBackoffJitter: DefaultMaxBackoffJitter,
			},
		},
	}
	for i := range testCases {
		tc := &testCases[i]
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			t.Logf("\nGIVEN: %s\nWHEN: %s\nTHEN: %s", tc.Given, tc.When, tc.Then)

			var o Options
			tc.o = &o

			if f := tc.Setup; f != nil {
				f(tc)
			}

			o.Sanitize()

			if f := tc.Teardown; f != nil {
				f(tc)
			}

			assert.Equal(t, tc.State, o)
		})
	}
}
