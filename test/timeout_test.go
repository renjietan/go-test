package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestWaitAndTimeout(t *testing.T) {
	var wg sync.WaitGroup
	done := make(chan struct{})

	// 启动工作 goroutine
	var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		longRunningTask(&ctx)
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("任务完成")
	case <-ctx.Done():
		fmt.Println("任务超时")
	}
}

func longRunningTask(ctx *context.Context) {
	// 模拟一个长时间运行的任务
	for i := 0; i < 10; i++ {
		select {
		case <-(*ctx).Done():
			// 收到取消信号，优雅退出
			fmt.Println("任务被取消")
			return
		case <-time.After(1 * time.Second):
			// 模拟工作
			fmt.Printf("执行第 %d 步\n", i+1)
		}
	}
	fmt.Println("任务完成")
}

func TestEasyTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务结束:", ctx.Err())
			return
		case <-time.After(1 * time.Second):
			fmt.Println("执行任务")
		}
	}
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		select {
	// 		case <-ctx.Done():
	// 			fmt.Println("任务超时")
	// 			return
	// 		case <-time.After(1 * time.Second):
	// 			fmt.Printf("执行第 %d 步\n", i+1)
	// 		}
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			fmt.Println("任务超时2")
	// 		default:
	// 			fmt.Println("==================")
	// 		}
	// 	}
	// }()
}
