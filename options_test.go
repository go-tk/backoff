package backoff_test

import (
	"testing"
	"time"

	. "github.com/go-tk/backoff"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestOptions_apply(t *testing.T) {
	type Workspace struct {
		O            Options
		B            Backoff
		ExpSt, ActSt string
	}
	tc := testcase.New().
		Step(0, func(t *testing.T, w *Workspace) {
			w.ExpSt = `
minDelay: 100ms
maxDelay: 1m40s
delayFactor: 2
maxDelayJitter: 1
maxNumberOfAttempts: 100
`[1:]
		}).
		Step(1, func(t *testing.T, w *Workspace) {
			w.O.Apply(&w.B)
			w.ActSt = w.B.DumpSettings()
		}).
		Step(2, func(t *testing.T, w *Workspace) {
			assert.Equal(t, w.ExpSt, w.ActSt)
		})
	testcase.RunListParallel(t,
		tc.Copy().
			Given("no option values").
			Then("settings should use default values"),
		tc.Copy().
			Given("invalid option values").
			Then("settings should use default values").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.O.MinDelay.Set(-1)
				w.O.MaxDelay.Set(-1)
				w.O.DelayFactor.Set(-1)
			}),
		tc.Copy().
			Given("valid option values").
			Then("settings should use these values").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.O.MinDelay.Set(1 * time.Second)
				w.O.MaxDelay.Set(2 * time.Second)
				w.O.DelayFactor.Set(3)
				w.O.MaxDelayJitter.Set(0.3)
				w.O.MaxNumberOfAttempts.Set(0)
			}).
			Step(1.5, func(t *testing.T, w *Workspace) {
				w.ExpSt = `
minDelay: 1s
maxDelay: 2s
delayFactor: 3
maxDelayJitter: 0.3
maxNumberOfAttempts: 0
`[1:]
			}),
		tc.Copy().
			Given("MinDelay option value > MaxDelay option value").
			Then("settings should use default values").
			Step(0.5, func(t *testing.T, w *Workspace) {
				w.O.MinDelay.Set(2 * time.Second)
				w.O.MaxDelay.Set(1 * time.Second)
			}),
	)
}
