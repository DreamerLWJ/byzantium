package utils

import (
	"context"
	"sync"
)

func GoWithWg(ctx context.Context, fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		fn()
	}()
}
