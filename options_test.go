package backoff_test

import (
	"strconv"
	"testing"
	"time"

	. "github.com/go-tk/backoff"
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
	type TestCase struct {
		Given, When, Then string
		Setup, Teardown   func(*TestCase)
		State             State

		t *testing.T
		o *Options
	}
	const (
		defaultMinDelay            time.Duration = 100 * time.Millisecond
		defaultMaxDelay            time.Duration = 100 * time.Second
		defaultDelayFactor         float64       = 2
		defaultMaxDelayJitter      float64       = 1
		defaultMaxNumberOfAttempts               = 100
	)
	testCases := []TestCase{
		{
			Then: "should overwrite options to default",
			State: State{
				MinDelay:            defaultMinDelay,
				MaxDelay:            defaultMaxDelay,
				DelayFactor:         defaultDelayFactor,
				MaxDelayJitter:      defaultMaxDelayJitter,
				MaxNumberOfAttempts: defaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MinDelay, MaxDelay, DelayFactor",
			Then:  "should overwrite options to default",
			Setup: func(tc *TestCase) {
				tc.o.MinDelay = -1
				tc.o.MaxDelay = -2
				tc.o.DelayFactor = -3
			},
			State: State{
				MinDelay:            defaultMinDelay,
				MaxDelay:            defaultMaxDelay,
				DelayFactor:         defaultDelayFactor,
				MaxDelayJitter:      defaultMaxDelayJitter,
				MaxNumberOfAttempts: defaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "options with specified values",
			Then:  "should not overwrite options",
			Setup: func(tc *TestCase) {
				tc.o.MinDelay = 1
				tc.o.MaxDelay = 2
				tc.o.DelayFactor = 3
				tc.o.MaxDelayJitter = 4
				tc.o.MaxNumberOfAttempts = 5
			},
			State: State{
				MinDelay:            1,
				MaxDelay:            2,
				DelayFactor:         3,
				MaxDelayJitter:      4,
				MaxNumberOfAttempts: 5,
			},
		},
		{
			Given: "option MaxDelay less than option MinDelay",
			Then:  "should overwrite option MaxDelay to MinDelay",
			Setup: func(tc *TestCase) {
				tc.o.MinDelay = 2
				tc.o.MaxDelay = 1
			},
			State: State{
				MinDelay:            2,
				MaxDelay:            2,
				DelayFactor:         defaultDelayFactor,
				MaxDelayJitter:      defaultMaxDelayJitter,
				MaxNumberOfAttempts: defaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MaxDelayJitter",
			Then:  "should overwrite option MaxDelayJitter to zero",
			Setup: func(tc *TestCase) {
				tc.o.MaxDelayJitter = -100
			},
			State: State{
				MinDelay:            defaultMinDelay,
				MaxDelay:            defaultMaxDelay,
				DelayFactor:         defaultDelayFactor,
				MaxNumberOfAttempts: defaultMaxNumberOfAttempts,
			},
		},
		{
			Given: "negative option MaxNumberOfAttempts",
			Then:  "should overwrite option MaxNumberOfAttempts to zero",
			Setup: func(tc *TestCase) {
				tc.o.MaxNumberOfAttempts = -100
			},
			State: State{
				MinDelay:       defaultMinDelay,
				MaxDelay:       defaultMaxDelay,
				DelayFactor:    defaultDelayFactor,
				MaxDelayJitter: defaultMaxDelayJitter,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			t.Logf("\nGIVEN: %s\nWHEN: %s\nTHEN: %s", tc.Given, tc.When, tc.Then)
			tc.t = t

			var o Options
			tc.o = &o

			if f := tc.Setup; f != nil {
				f(&tc)
			}

			o.Sanitize()

			if f := tc.Teardown; f != nil {
				f(&tc)
			}

			var state State
			state.MinDelay = o.MinDelay
			state.MaxDelay = o.MaxDelay
			state.DelayFactor = o.DelayFactor
			state.MaxDelayJitter = o.MaxDelayJitter
			state.DelayFuncIsNil = o.DelayFunc == nil
			state.MaxNumberOfAttempts = o.MaxNumberOfAttempts
			assert.Equal(t, tc.State, state)
		})
	}
}
