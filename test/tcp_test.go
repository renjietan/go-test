package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TcpClients struct {
}

type Action = string

const (
	Connect    Action = "connect"
	Login      Action = "login"
	HeartBeat  Action = "heartbeat"
	DisConnect Action = "idsconnect"
)

type TcpClient struct {
	Conn   net.Conn
	Action chan Action
	Runing atomic.Bool
	Mu     sync.Mutex
}

func (t *TcpClient) Connect(addr string) {
	conn, err := net.Dial("tcp", addr)
	t.Conn = conn
	if err != nil {
		fmt.Println("Connect err:", err)
	}
	t.Action = make(chan Action, 10)
	t.Runing.Store(true)
	go t.OnMessage()
	go t.sendMessage("hello")
	go t.HandleData()

}

func (t *TcpClient) OnMessage() {
	defer t.Disconnect()
	reader := bufio.NewReader(t.Conn)
	for t.Runing.Load() {
		t.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		message, err := reader.ReadString('o')
		fmt.Println("message:", message)
		fmt.Println("err:", err)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF: 服务器关闭连接")
			} else {
				fmt.Printf("读取消息失败: %v\n", err)
			}
			return
		}
		t.Action <- HeartBeat
		fmt.Println("message:", message)
	}
}

// sendMessage 发送消息
func (t *TcpClient) sendMessage(message string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t.Conn == nil {
		return
	}

	// message += "\n"
	_, err := t.Conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("发送消息失败: %v\n", err)
		t.Disconnect()
		return
	}
	fmt.Printf("发送消息: %s\n", strings.TrimSpace(message))

	// // 启动超时检查
	// pool.Submit(func() {
	// 	select {
	// 	case <-time.After(5 * time.Second):
	// 		fmt.Println("发送消息后5秒内未收到回复，断开连接")
	// 	case <-c.timeoutCh:
	// 		// 收到消息，取消超时
	// 	}
	// })
}

func (t *TcpClient) HandleData() {
	// defer t.Disconnect()
	for t.Runing.Load() {
		select {
		case action := <-t.Action:
			go t.sendMessage("hello")
			fmt.Println("action:", action)
		default:

		}
	}
}

func (t *TcpClient) Disconnect() {
	fmt.Println("执行  断开连接  方法")
	if !t.Runing.Load() {
		return
	}
	t.Runing.Store(false)

	if t.Conn != nil {
		t.Conn.Close()
		t.Conn = nil
	}
}

func Test_tcp(t *testing.T) {
	tcp := &TcpClient{}
	defer tcp.Disconnect()
	tcp.Connect("8.135.10.183:30558")
	time.Sleep(10000 * time.Second)
}
