package main

import (
	"fmt"
	"gframe-ants-tcp/pool"
)

func main() {
	// 初始化协程池
	err := pool.InitPool(100)
	if err != nil {
		panic(fmt.Sprintf("初始化协程池失败: %v", err))
	}
	defer pool.Release()

	// 启动 TCP 服务器
	if err := StartTCPServer(); err != nil {
		panic(fmt.Sprintf("启动TCP服务器失败: %v", err))
	}
}
