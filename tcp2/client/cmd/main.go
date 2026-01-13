package main

import (
	"fmt"
	"gframe-ants-tcp/pool"
	"gframe-ants-tcp/tcp2/client"
	"sync"
	"time"
)

func main() {
	// 初始化协程池
	err := pool.InitPool(100)
	if err != nil {
		panic(fmt.Sprintf("初始化协程池失败: %v", err))
	}
	defer pool.Release()

	serverAddr := "127.0.0.1:8080"
	clientCount := 10

	fmt.Printf("开始创建 %d 个TCP客户端连接到 %s\n", clientCount, serverAddr)

	var wg sync.WaitGroup

	// 循环创建10个TCP客户端
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			createAndStartClient(index, serverAddr)
		}(i)

		// 稍微延迟，避免同时连接造成压力
		time.Sleep(100 * time.Millisecond)
	}

	// 等待所有客户端创建完成
	wg.Wait()

	fmt.Printf("\n所有客户端创建完成，当前客户端数: %d, 已连接数: %d\n",
		client.GlobalManager.GetClientCount(),
		client.GlobalManager.GetConnectedCount())

	// 保持程序运行
	select {}
}

// createAndStartClient 创建并启动客户端
func createAndStartClient(index int, serverAddr string) {
	// 在管理器中添加客户端信息
	clientInfo := client.GlobalManager.AddClient(serverAddr)
	fmt.Printf("创建客户端 [ID: %d] 准备连接到 %s\n", clientInfo.ID, serverAddr)

	// 创建TCP客户端实例
	tcpClient := client.NewTCPClient(clientInfo.ID, serverAddr)

	// 保存客户端实例，便于外部访问
	client.StoreClient(clientInfo.ID, tcpClient)

	// 连接到服务端
	err := tcpClient.Connect()
	if err != nil {
		fmt.Printf("客户端 [ID: %d] 连接失败: %v\n", clientInfo.ID, err)
		client.GlobalManager.RemoveClient(clientInfo.ID)
		client.RemoveClientInstance(clientInfo.ID)
		return
	}

	fmt.Printf("客户端 [ID: %d] 启动成功\n", clientInfo.ID)
}
