package main

import (
	"fmt"
	"gframe-ants-tcp/pool"
	"time"
)

// 示例：在其他文件中使用协程池
func exampleUsePool() {
	// 直接使用 pool.Submit 提交任务，非常方便
	for i := 0; i < 10; i++ {
		taskID := i
		pool.Submit(func() {
			fmt.Printf("任务 %d 正在协程池中执行\n", taskID)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("任务 %d 执行完成\n", taskID)
		})
	}

	// 可以查看协程池状态
	fmt.Printf("协程池状态 - 运行中: %d, 空闲: %d, 容量: %d\n",
		pool.Running(), pool.Free(), pool.Cap())
}
