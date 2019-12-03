package util

import (
	"context"
	"time"
)

func Watch(ctx context.Context, notify chan<- bool) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
	stop:
		for {
			select {
			case <-ctx.Done():
				Info("terminating: context cancelled")
				notify <- true
				break stop
			case <-ticker.C:
				if ctx.Err() != nil {
					Info("terminating: context cancelled")
					notify <- true
					break stop
				}
			}
		}
		ticker.Stop()
	}()
}
