package main

import (
	"fmt"
	"gframe-ants-tcp/tcp2/client"
)

// 示例：从程序任意一处访问客户端管理器
// 只需导入 client 包，即可使用 client.GlobalManager 访问所有客户端

func AccessClientsFromAnywhere() {
	// 获取所有客户端
	clients := client.GlobalManager.GetAllClients()
	fmt.Printf("当前连接的客户端数量: %d\n", len(clients))

	// 获取客户端总数
	count := client.GlobalManager.GetClientCount()
	fmt.Printf("当前活跃客户端数: %d\n", count)

	// 根据ID获取特定客户端
	if clientInfo, ok := client.GlobalManager.GetClient(1); ok {
		fmt.Printf("找到客户端 ID: %d, 地址: %s\n", clientInfo.ID, clientInfo.Addr)
		// 可以向该客户端发送消息
		clientInfo.Conn.Write([]byte("来自任意位置的消息\n"))
	}

	// 广播消息给所有客户端
	client.GlobalManager.Broadcast("广播消息: 来自程序任意一处\n")
}
