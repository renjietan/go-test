package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 定时器 包含 停止 功能
func TestTicker(t *testing.T) {
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		t := time.NewTicker(time.Second)
		defer func() {
			t.Stop()
			wg.Done()
		}()
		for {
			select {
			case <-t.C:
				fmt.Println("定时器触发")
			case <-done:
				fmt.Println("定时器停止")
				return
			}
		}
	}()
	defer close(done)
	go func() {
		time.Sleep(3 * time.Second)
		wg.Done()
		done <- struct{}{}
	}()
	wg.Wait()
}
