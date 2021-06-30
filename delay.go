package backoff

import (
	"math/rand"
	"time"
)

type delay struct {
	d             time.Duration
	isInitialized bool
	rand          *rand.Rand
}

func (d *delay) Update(min, max time.Duration, factor, maxJitter float64) time.Duration {
	if d.isInitialized {
		d.d = time.Duration(float64(d.d) * factor)
		if d.d > max {
			d.d = max
		}
	} else {
		d.d = min
		d.isInitialized = true
	}
	return d.withJitter(maxJitter)
}

func (d *delay) withJitter(maxJitter float64) time.Duration {
	if maxJitter <= 0 {
		return d.d
	}
	if d.rand == nil {
		d.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	jitter := maxJitter * (2*d.rand.Float64() - 1)
	dd := d.d + time.Duration(float64(d.d)*jitter)
	return dd
}
