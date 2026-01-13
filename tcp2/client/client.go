package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
	"gframe-ants-tcp/pool"
)

// TCPClient TCP客户端结构
type TCPClient struct {
	ID          int64
	Conn        net.Conn
	Addr        string
	IsActive    bool
	IsConnected bool
	stopChan    chan struct{}
	replyChan   chan struct{} // 用于通知收到服务端回复
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// NewTCPClient 创建新的TCP客户端
func NewTCPClient(id int64, addr string) *TCPClient {
	return &TCPClient{
		ID:          id,
		Addr:        addr,
		IsActive:    true,
		IsConnected: false,
		stopChan:    make(chan struct{}),
		replyChan:   make(chan struct{}, 1), // 带缓冲，避免阻塞
	}
}

// Connect 连接到服务端
func (c *TCPClient) Connect() error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("连接服务端失败: %v", err)
	}

	c.mu.Lock()
	c.Conn = conn
	c.IsConnected = true
	c.mu.Unlock()

	// 更新管理器中的状态
	GlobalManager.UpdateClientStatus(c.ID, true)

	fmt.Printf("客户端 [ID: %d] 成功连接到服务端: %s\n", c.ID, c.Addr)

	// 启动消息监听
	c.wg.Add(1)
	pool.Submit(func() {
		defer c.wg.Done()
		c.listenMessages()
	})

	// 启动心跳
	c.wg.Add(1)
	pool.Submit(func() {
		defer c.wg.Done()
		c.startHeartbeat()
	})

	return nil
}

// listenMessages 监听服务端消息
func (c *TCPClient) listenMessages() {
	defer func() {
		c.mu.Lock()
		if c.Conn != nil {
			c.Conn.Close()
		}
		c.IsConnected = false
		c.IsActive = false
		c.mu.Unlock()
		GlobalManager.UpdateClientStatus(c.ID, false)
		fmt.Printf("客户端 [ID: %d] 断开连接\n", c.ID)
	}()

	reader := bufio.NewReader(c.Conn)
	for {
		select {
		case <-c.stopChan:
			return
		default:
			// 设置读取超时
			c.mu.RLock()
			conn := c.Conn
			c.mu.RUnlock()

			if conn == nil {
				return
			}

			// 设置读取超时
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))

			message, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Printf("客户端 [ID: %d] 服务端关闭连接\n", c.ID)
				} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// 超时继续循环
					continue
				} else {
					fmt.Printf("客户端 [ID: %d] 读取消息失败: %v\n", c.ID, err)
				}
				return
			}

			// 处理接收到的消息
			fmt.Printf("客户端 [ID: %d] 收到服务端消息: %s", c.ID, string(message))
			
			// 通知心跳协程收到服务端回复
			select {
			case c.replyChan <- struct{}{}:
			default:
				// 如果channel已满，跳过（避免阻塞）
			}
		}
	}
}

// startHeartbeat 启动心跳功能
// 心跳逻辑：连接后立即发送一次心跳，等待服务端回复后，延迟3秒再发送一次心跳
func (c *TCPClient) startHeartbeat() {
	// 连接后立即发送一次心跳
	c.sendHeartbeat()

	// 等待服务端回复
	for {
		select {
		case <-c.stopChan:
			return
		case <-c.replyChan:
			// 收到服务端回复后，延迟3秒再发送一次心跳
			time.Sleep(3 * time.Second)
			
			c.mu.RLock()
			isConnected := c.IsConnected
			c.mu.RUnlock()

			if !isConnected {
				return
			}

			c.sendHeartbeat()
		}
	}
}

// sendHeartbeat 发送心跳消息
func (c *TCPClient) sendHeartbeat() error {
	c.mu.RLock()
	conn := c.Conn
	isConnected := c.IsConnected
	c.mu.RUnlock()

	if !isConnected || conn == nil {
		return fmt.Errorf("客户端未连接")
	}

	heartbeatMsg := fmt.Sprintf("heartbeat from client %d: %s\n", c.ID, time.Now().Format("2006-01-02 15:04:05"))
	_, err := conn.Write([]byte(heartbeatMsg))
	if err != nil {
		fmt.Printf("客户端 [ID: %d] 发送心跳失败: %v\n", c.ID, err)
		return err
	}

	fmt.Printf("客户端 [ID: %d] 发送心跳: %s", c.ID, heartbeatMsg)
	return nil
}

// Close 关闭客户端连接
func (c *TCPClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsActive {
		return
	}

	c.IsActive = false
	c.IsConnected = false
	close(c.stopChan)

	if c.Conn != nil {
		c.Conn.Close()
	}

	// 等待所有协程退出
	c.wg.Wait()

	GlobalManager.UpdateClientStatus(c.ID, false)
	fmt.Printf("客户端 [ID: %d] 已关闭\n", c.ID)
}

// SendMessage 发送消息到服务端
func (c *TCPClient) SendMessage(message string) error {
	c.mu.RLock()
	conn := c.Conn
	isConnected := c.IsConnected
	c.mu.RUnlock()

	if !isConnected || conn == nil {
		return fmt.Errorf("客户端未连接")
	}

	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Printf("客户端 [ID: %d] 发送消息失败: %v\n", c.ID, err)
		return err
	}

	return nil
}
