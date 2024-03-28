package background

import (
	"context"
	"time"
)

func Ticker(ctx context.Context, every time.Duration, async bool, f func()) {
	last := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Since(last) >= every {
				if async {
					go f()
				} else {
					f()
				}
				last = time.Now()
			}
		}
	}
}
