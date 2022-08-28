package backoff

import (
	"math/rand"
	"time"

	"github.com/benbjohnson/clock"
)

type delay struct {
	d    time.Duration
	rand *rand.Rand
}

func (d *delay) Update(min, max time.Duration, factor, maxJitter float64, clock clock.Clock) time.Duration {
	if d.d < min {
		d.d = min
	} else {
		d.d = time.Duration(float64(d.d) * factor)
		if d.d > max {
			d.d = max
		}
	}
	return d.withJitter(maxJitter, clock)
}

func (d *delay) withJitter(maxJitter float64, clock clock.Clock) time.Duration {
	if maxJitter <= 0 {
		return d.d
	}
	if d.rand == nil {
		d.rand = rand.New(rand.NewSource(clock.Now().UnixNano()))
	}
	jitter := maxJitter * (2*d.rand.Float64() - 1)
	dd := d.d + time.Duration(float64(d.d)*jitter)
	return dd
}
