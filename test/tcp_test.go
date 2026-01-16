package main

import (
	"bufio"
	"context"
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
	Client *TcpClient
}

type Action = string

const (
	Connect    Action = "connect"
	Login      Action = "login"
	HeartBeat  Action = "heartbeat"
	DisConnect Action = "disconnect"
)

func NewTcpClient(addr string) {
	tcp := &TcpClient{}
	defer tcp.Disconnect("程序结束")
	tcp.Connect("8.135.10.183:30558")
}

type TcpClient struct {
	Conn            net.Conn
	Action          chan Action
	Running         atomic.Bool
	Mu              sync.Mutex
	CleanTimeoutTag chan struct{}
	Once            sync.Once
	Done            chan struct{}
	timer           *time.Timer
}

func (t *TcpClient) Connect(addr string) {
	defer func() {
		t.timer = time.NewTimer(10 * time.Second)
		t.Action <- HeartBeat
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connect err:", err)
		return
	}
	t.Conn = conn
	// 注意 此处一定要 初始化
	t.Action = make(chan Action, 100)
	t.CleanTimeoutTag = make(chan struct{}, 100)
	t.Done = make(chan struct{})
	t.Running.Store(true)
	go t.OnMessage()
	go t.HandleData()
}

func (t *TcpClient) OnMessage() {
	defer t.Disconnect("onMessage 循环接收消息 defer ")
	reader := bufio.NewReader(t.Conn)
	for t.Running.Load() {
		t.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		message, err := reader.ReadString('o')
		fmt.Println("收到消息:", message)
		fmt.Println("err:", err)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF: 服务器关闭连接")
			} else {
				fmt.Printf("读取消息失败: %v\n", err)
			}
			return
		}
		select {
		case t.CleanTimeoutTag <- struct{}{}:
		case <-t.Done:
			return
		}
		time.Sleep(3 * time.Second)
		t.Action <- HeartBeat
	}
}

// 发送消息
func (t *TcpClient) sendMessage(message string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t.Conn == nil {
		return
	}
	_, err := t.Conn.Write([]byte(message))
	if err != nil {
		t.Disconnect(fmt.Sprintf("发送消息失败: %v\n", err))
		return
	}
	fmt.Printf("发送消息: %s\n", strings.TrimSpace(message))
	go func() {
		t.Mu.Lock()
		t.timer.Reset(10 * time.Second)
		defer t.Mu.Unlock()
		select {
		case <-t.timer.C:
			t.Disconnect("消息接收超时")
		case <-t.CleanTimeoutTag:
			fmt.Println("清除延时器")
			t.timer.Stop()
		}

	}()
}

func (t *TcpClient) HandleData() {
	defer t.Disconnect("消息处理后，开始发送")
	for t.Running.Load() {
		select {
		case action := <-t.Action:
			go t.sendMessage("hello")
			fmt.Println("action:", action)
		case <-t.Done:
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *TcpClient) Disconnect(log_str string) {
	t.Once.Do(func() {
		fmt.Errorf("=================== %v ====================\n", log_str)
		t.Running.Store(false)
		close(t.Done) // 通知所有 goroutine 退出
		if t.Conn != nil {
			t.Conn.Close()
			t.Conn = nil
		}
	})
}

func Test_tcp(t *testing.T) {
	tcp := &TcpClient{}
	defer tcp.Disconnect("程序结束")

	tcp.Connect("8.135.10.183:30558")

	// 用 context 超时控制测试，改为10秒
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	select {
	case <-tcp.Done:
		fmt.Println("连接已断开")
	case <-ctx.Done():
		fmt.Println("测试超时")
	}
}
