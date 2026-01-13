package client

import "fmt"

// 示例：展示如何在程序任意一处访问客户端

// ExampleAccessClients 示例：访问客户端管理器
func ExampleAccessClients() {
	// 1. 获取所有客户端信息
	clients := GlobalManager.GetAllClients()
	fmt.Printf("当前客户端数量: %d\n", len(clients))

	// 2. 获取客户端实例
	allInstances := GetAllClientInstances()
	fmt.Printf("当前客户端实例数: %d\n", len(allInstances))

	// 3. 根据ID获取特定客户端实例
	if clientInstance, ok := GetClientInstance(1); ok {
		fmt.Printf("找到客户端实例 ID: %d, 地址: %s\n", clientInstance.ID, clientInstance.Addr)
		
		// 可以向该客户端发送消息
		clientInstance.SendMessage("来自外部的消息")
	}

	// 4. 获取客户端统计信息
	fmt.Printf("客户端总数: %d, 已连接数: %d\n",
		GlobalManager.GetClientCount(),
		GlobalManager.GetConnectedCount())
}
