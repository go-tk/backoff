package backoff

import (
	"time"
)

type timer struct {
	t *time.Timer
}

func (t *timer) Start(timeout time.Duration) <-chan struct{} {
	timedOut := make(chan struct{})
	t.t = time.AfterFunc(timeout, func() { close(timedOut) })
	return timedOut
}

func (t *timer) Stop() { t.t.Stop() }
