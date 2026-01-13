package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	tcp_client_enum "gframe-ants-tcp/tcp/client/enum"
	"net"
	"time"
)

type Client struct {
	Object net.Conn
	Key    int
	Action chan tcp_client_enum.ClientActionEnum
}

var pools = make(chan Client, 40)

func createClient(key int) {
	// 建立TCP连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic("创建客户端失败" + err.Error())
	}
	// 注意：不要在这里 defer conn.Close()，因为连接需要持续使用
	// 连接会在客户端退出时或出错时关闭

	// 初始化 Action channel（带缓冲，避免阻塞）
	actionChan := make(chan tcp_client_enum.ClientActionEnum, 10)

	client := Client{
		Object: conn,
		Key:    key,
		Action: actionChan,
	}

	// 启动监听服务端消息的 goroutine
	go listenServerMessages(conn)

	// 启动监听 Action channel 的 goroutine
	go handleClientActions(client)

	go heartBeat(client.Object)
	pools <- client
}

func heartBeat(conn net.Conn) {
	t := time.NewTicker(3 * time.Second)
	for {
		time.Sleep(2 * time.Second)
		select {
		case <-t.C:
			now := time.Now().Format("2006-01-02 15:04:05")
			_, err := conn.Write([]byte(now + " : " + "\n"))
			if err != nil {
				fmt.Println("Error writing to connection:", err)
			}
		default:
			fmt.Println("======================")
		}
	}
}

func listenServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// 读取直到遇到换行符
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("连接读取错误: %v\n", err)
			break // 连接已关闭或出错
		}
		// 处理接收到的消息（这里简单打印）
		fmt.Println("收到服务端消息:", string(message))
	}
}

// handleClientActions 通过 select 循环读取 Action 字段并处理
func handleClientActions(client Client) {
	// 使用 for range 遍历 channel，当 channel 关闭时自动退出
	for action := range client.Action {
		switch action {
		case tcp_client_enum.ClientActionConnect:
			fmt.Printf("客户端 %d 执行连接动作\n", client.Key)
			// 可以在这里执行连接相关的操作

		case tcp_client_enum.ClientActionLogin:
			fmt.Printf("客户端 %d 执行登录动作\n", client.Key)
			// 发送登录消息
			data, _ := json.Marshal(map[string]interface{}{
				"Name":   "张三",
				"Age":    "age",
				"Key":    client.Key,
				"Action": "login",
			})
			client.Object.Write(data)

		case tcp_client_enum.ClientActionHeartbeat:
			fmt.Printf("客户端 %d 执行心跳动作\n", client.Key)
			// 发送心跳消息
			now := time.Now().Format("2006-01-02 15:04:05")
			heartbeatMsg := fmt.Sprintf("heartbeat: %s\n", now)
			client.Object.Write([]byte(heartbeatMsg))
		}
	}
	fmt.Printf("客户端 %d 的 Action channel 已关闭\n", client.Key)
}

// 高并发客户端示例
func main() {
	for i := 0; i < 4; i++ {
		go func(index int) {
			createClient(index)
			// data, _ := json.Marshal(map[string]interface{}{"Name": "张三", "Age": "age", "Key": client.Key})
			// client.Object.Write(data)
			// defer conn.Close()
		}(i)
	}

	go func() {
		for {
			select {
			case obj := <-pools:
				data, _ := json.Marshal(map[string]interface{}{"Name": "张三", "Age": "age", "Key": obj.Key})
				obj.Object.Write(data)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// // 等待所有goroutine完成
	// for i := 0; i < 40; i++ {
	// 	sem <- struct{}{}
	// }
	time.Sleep(5000 * time.Second)
}
