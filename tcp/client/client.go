package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"gframe-ants-tcp/pool"
)

// Client TCP客户端结构体
type Client struct {
	conn       net.Conn
	serverAddr string
	mu         sync.Mutex
	running    bool
	timeoutCh  chan bool // 超时通道
}

// NewClient 创建新的TCP客户端
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		timeoutCh:  make(chan bool, 1),
	}
}

// Connect 连接到服务器
func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.serverAddr)
	if err != nil {
		return err
	}
	fmt.Printf("连接到服务器: %s\n", c.serverAddr)

	c.running = true
	// 连接后立即发送 "hello-1"
	pool.Submit(func() {
		c.sendMessage("hello-1")
	})

	// 启动监听协程
	pool.Submit(func() {
		c.onMessage()
	})

	return nil
}

// sendMessage 发送消息
func (c *Client) sendMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return
	}

	// message += "\n"
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("发送消息失败: %v\n", err)
		c.disconnect()
		return
	}
	fmt.Printf("发送消息: %s\n", strings.TrimSpace(message))

	// 启动超时检查
	pool.Submit(func() {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("发送消息后5秒内未收到回复，断开连接")
		case <-c.timeoutCh:
			// 收到消息，取消超时
		}
	})
}

// onMessage 监听服务器消息
func (c *Client) onMessage() {
	defer c.disconnect()
	reader := bufio.NewReader(c.conn)
	fmt.Println("running:", c.running)
	for c.running {
		// 设置读取超时
		// c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		message, err := reader.ReadString('\n')
		fmt.Println("msg==================", message)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF: 服务器关闭连接")
			} else {
				fmt.Printf("读取消息失败: %v\n", err)
			}
			return
		}

		// message = strings.TrimSpace(string(message))
		fmt.Printf("收到服务器消息: %s\n", strings.TrimSpace(string(message)))

		select {
		case c.timeoutCh <- true:
		default:
		}

		// 收到消息后，3秒后回复 "hello-3"
		pool.Submit(func() {
			time.Sleep(3 * time.Second)
			c.sendMessage("hello-3")
		})
	}
}

func (c *Client) disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}
	c.running = false

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		fmt.Println("客户端断开连接(disconnected)")
	}
}

// IsConnected 检查是否连接
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running && c.conn != nil
}

// Close 关闭客户端连接
func (c *Client) Close() {
	c.disconnect()
}
