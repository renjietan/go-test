package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {
	var wg sync.WaitGroup
	var b atomic.Bool
	wg.Add(1001)
	for v := range 1000 {
		go func(index int, data *atomic.Bool) {
			defer func() {
				wg.Done()
			}()
			time.Sleep(time.Second)
			data.Store(index > 500)
			fmt.Printf("v: %v  %v\n", index, data.Load())
		}(v, &b)
	}
	go func(data *atomic.Bool) {
		defer wg.Done()
		for b.Load() {
			fmt.Println("=========== for ============")
		}
	}(&b)
	wg.Wait()
}
