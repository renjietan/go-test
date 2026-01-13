package main

import (
	"bufio"
	"fmt"
	"gframe-ants-tcp/pool"
	"gframe-ants-tcp/tcp2/client"
	"io"
	"net"
)

// StartTCPServer 启动 TCP 服务器
func StartTCPServer() error {
	// 监听地址 127.0.0.1:8080
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return fmt.Errorf("创建TCP监听失败: %v", err)
	}
	defer listener.Close()

	fmt.Println("TCP服务器已启动，监听地址: 127.0.0.1:8080")

	// 持续接受客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}

		// 添加客户端到管理器
		clientInfo := client.GlobalManager.AddClient(conn)
		fmt.Printf("新客户端连接 [ID: %d, 地址: %s]，当前客户端数: %d\n",
			clientInfo.ID, clientInfo.Addr, client.GlobalManager.GetClientCount())

		// 使用协程池处理客户端连接
		pool.Submit(func() {
			handleClient(clientInfo)
		})
	}
}

// handleClient 处理客户端连接和消息
func handleClient(clientInfo *client.Info) {
	defer func() {
		clientInfo.Conn.Close()
		client.GlobalManager.RemoveClient(clientInfo.ID)
		fmt.Printf("客户端 [ID: %d, 地址: %s] 断开连接，当前客户端数: %d\n",
			clientInfo.ID, clientInfo.Addr, client.GlobalManager.GetClientCount())
	}()

	reader := bufio.NewReader(clientInfo.Conn)

	// 持续监听客户端消息
	for {
		// 读取消息（按行读取）
		message, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端 [ID: %d] 正常断开连接\n", clientInfo.ID)
			} else {
				fmt.Printf("读取客户端 [ID: %d] 消息失败: %v\n", clientInfo.ID, err)
			}
			return
		}

		// 打印收到的消息
		fmt.Printf("收到客户端 [ID: %d, 地址: %s] 消息: %s",
			clientInfo.ID, clientInfo.Addr, string(message))

		// 固定回复 "hello client"
		response := "hello client\n"
		_, err = clientInfo.Conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("向客户端 [ID: %d] 发送回复失败: %v\n", clientInfo.ID, err)
			return
		}
	}
}
