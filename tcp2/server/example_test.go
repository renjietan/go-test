package main

import (
	"fmt"
	"gframe-ants-tcp/tcp2/client"
	"testing"
	"time"
)

// TestAccessClients 测试：展示如何在程序任意一处访问客户端管理器
// 只需导入 "gframe-ants-tcp/tcp2/client" 包即可使用 client.GlobalManager
func TestAccessClients(t *testing.T) {
	// 1. 获取所有客户端
	clients := client.GlobalManager.GetAllClients()
	fmt.Printf("当前连接的客户端数量: %d\n", len(clients))

	// 2. 遍历所有客户端
	for _, clientInfo := range clients {
		fmt.Printf("客户端 ID: %d, 地址: %s, 活跃: %v\n",
			clientInfo.ID, clientInfo.Addr, clientInfo.IsActive)
	}

	// 3. 根据ID获取特定客户端
	if clientInfo, ok := client.GlobalManager.GetClient(1); ok {
		fmt.Printf("找到客户端 ID: %d, 地址: %s\n", clientInfo.ID, clientInfo.Addr)
		// 可以向该客户端发送消息
		clientInfo.Conn.Write([]byte("来自服务器的消息\n"))
	}

	// 4. 获取客户端总数
	count := client.GlobalManager.GetClientCount()
	fmt.Printf("当前活跃客户端数: %d\n", count)

	// 5. 广播消息给所有客户端
	client.GlobalManager.Broadcast("广播消息: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
}

// Example 包级别示例：展示如何在程序任意一处访问客户端管理器
func Example() {
	// 获取所有客户端
	clients := client.GlobalManager.GetAllClients()
	fmt.Printf("当前连接的客户端数量: %d\n", len(clients))

	// 获取客户端总数
	count := client.GlobalManager.GetClientCount()
	fmt.Printf("当前活跃客户端数: %d\n", count)

	// 广播消息给所有客户端
	client.GlobalManager.Broadcast("广播消息\n")
}
