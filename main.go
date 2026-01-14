package main

import (
	"fmt"
	"time"

	"gframe-ants-tcp/pool"
	"gframe-ants-tcp/tcp/client"
)

func main() {
	// 初始化协程池，设置大小为 100
	err := pool.InitPool(100)
	if err != nil {
		panic(fmt.Sprintf("初始化协程池失败: %v", err))
	}
	defer pool.Release() // 程序退出时释放协程池

	fmt.Printf("协程池初始化成功，容量: %d\n", pool.Cap())

	// 初始化TCP服务器
	// tcpServer := server.NewServer()
	// go func() {
	// 	err := tcpServer.Start(":8080")
	// 	if err != nil {
	// 		fmt.Printf("启动TCP服务器失败: %v\n", err)
	// 	}
	// }()
	// defer tcpServer.Stop()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 初始化TCP客户端
	tcpClient := client.NewClient("8.135.10.183:30558")
	err = tcpClient.Connect()
	if err != nil {
		fmt.Printf("连接到服务器失败: %v\n", err)
		return
	}
	defer tcpClient.Close()
	time.Sleep(1000 * time.Second)
}
