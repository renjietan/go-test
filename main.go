package main

import (
	"fmt"
	"gframe-ants-tcp/pool"
)

func main() {
	// 初始化协程池，设置大小为 100
	err := pool.InitPool(100)
	if err != nil {
		panic(fmt.Sprintf("初始化协程池失败: %v", err))
	}
	defer pool.Release() // 程序退出时释放协程池

	fmt.Printf("协程池初始化成功，容量: %d\n", pool.Cap())

	// 示例：使用协程池执行任务
	for index, item := range []int{1, 2, 3, 4, 5, 6} {
		if item == 4 {
			continue
		}
		// 使用协程池提交任务
		pool.Submit(func() {
			fmt.Printf("第%d个元素的值是: %d\n", index, item)
		})
	}
}
