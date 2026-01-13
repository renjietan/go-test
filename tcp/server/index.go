package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

type Client struct {
	Object net.Conn
	Key    int
}

var clients sync.Map

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("创建失败", err)
		return
	}
	defer listener.Close()
	fmt.Println("服务器已启动，等待连接...")
	// go heartBeat()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败", err)
			continue
		}
		clients.Delete(conn)
		clients.Store(conn, 0)
		go handleConnection(conn)
	}

}

// func heartBeat() {
// 	t := time.NewTicker(3 * time.Second)
// 	for {
// 		time.Sleep(2 * time.Second)
// 		select {
// 		case <-t.C:
// 			now := time.Now().Format("2006-01-02 15:04:05")
// 			clients.Range(func(key, value any) bool {
// 				fmt.Println("value:", value)
// 				conn := key.(net.Conn)
// 				_, err := conn.Write([]byte(now + " : " + string(value.(int)) + "\n"))
// 				if err != nil {
// 					fmt.Println("Error writing to connection:", err)
// 				}
// 				return true
// 			})
// 		default:
// 			fmt.Println("======================")
// 		}
// 	}
// }

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("新客户端连接: %s\n", conn.RemoteAddr().String())

	// 使用 bufio.Reader 来读取数据，支持按行读取
	reader := bufio.NewReader(conn)

	// 持续读取客户端消息
	for {
		// 读取直到遇到换行符，或者读取固定长度的数据
		message, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端 %s 断开连接\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			// 从客户端列表中移除
			clients.Delete(conn)
			return
		}

		// 处理接收到的消息
		fmt.Printf("收到消息 [%s]: %s", conn.RemoteAddr().String(), string(message))

		// 发送响应
		response := "Message received successfully\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			clients.Delete(conn)
			return
		}
	}
}
