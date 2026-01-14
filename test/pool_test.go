package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gframe-ants-tcp/pool"
)

var wg sync.WaitGroup

func Test_tt1(t *testing.T) {
	err := pool.InitPool(100)
	if err != nil {
		panic(fmt.Sprintf("初始化协程池失败: %v", err))
	}
	defer pool.Release()

	for i := 0; i < 4; i++ {
		idx := i // 捕获循环变量
		pool.Submit(func() {
			ticker := time.NewTicker(time.Duration(idx+1) * time.Second)
			defer ticker.Stop()
			now := time.Now()
			fmt.Printf("协程 %d 启动，当前时间: %s，定时器间隔: %d秒\n", idx, now.Format("2006-01-02 15:04:05"), idx+1)
			// 持续等待定时器触发
			count := 0
			for v := range ticker.C {
				count++
				// 计算与now 的时间差
				duration := v.Sub(now)
				duration_num := duration.Milliseconds()
				fmt.Printf("时间: %v 协程 %d 第%d次触发，耗时: %v\n", v.Format("2006-01-02 15:04:05"), idx, count, duration)
				if duration_num > 10000 {
					fmt.Println("任务执行时间过长")
					ticker.Stop()
					break
				}
				// 可以设置触发次数限制，比如只触发5次
				if count >= 5 {
					fmt.Printf("协程 %d 已完成5次触发，退出\n", idx)
					break
				}
			}
		})
	}
	// time.Sleep(100 * time.Second) // 等待足够的时间让所有定时器触发
}

func Test_chan(t *testing.T) {
	cn := make(chan int, 10)
	wg.Add(10) // 添加等待计数，对应10个goroutine

	// 启动一个goroutine接收数据
	go func() {
		// for v := range cn {
		// 	fmt.Println("v:", v)
		// }
		for {
			select {
			case v, _ := <-cn:
				fmt.Println("v:", v)
			case t := <-time.After(3 * time.Second):
				fmt.Println("超时：", t.Format("2006-01-02 15:04:05"))
			}
		}
	}()

	// 启动10个goroutine发送数据
	for v := range 10 {
		go func(index int) {
			defer wg.Done()
			cn <- index
			fmt.Printf("已发送: %d\n", index)
			time.Sleep(time.Duration(index) * time.Second)
		}(v)
	}

	for i := 10; i < 20; i++ {
		cn <- i
	}

	wg.Wait()
	close(cn)
	time.Sleep(100 * time.Second)
}
