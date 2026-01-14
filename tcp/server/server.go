package server

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

// Server TCP服务器结构体
type Server struct {
	listener net.Listener
	mu       sync.Mutex
	conns    map[net.Conn]bool // 跟踪连接
}

// NewServer 创建新的TCP服务器
func NewServer() *Server {
	return &Server{
		conns: make(map[net.Conn]bool),
	}
}

func (s *Server) Start(address string) error {
	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("TCP服务器启动，监听地址: %s\n", address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("接受连接失败: %v\n", err)
			continue
		}
		fmt.Printf("新客户端连接: %s\n", conn.RemoteAddr().String())

		// 添加到连接列表
		s.mu.Lock()
		s.conns[conn] = true
		s.mu.Unlock()

		// 使用协程池处理连接
		pool.Submit(func() {
			s.handleConnection(conn)
		})
	}
}

// handleConnection 处理单个连接
func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		// 移除连接
		s.mu.Lock()
		delete(s.conns, conn)
		s.mu.Unlock()
		conn.Close()
		fmt.Printf("客户端断开连接: %s\n", conn.RemoteAddr().String())
	}()

	reader := bufio.NewReader(conn)
	for {
		// 设置读取超时
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("客户端关闭连接: %s\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("读取消息失败: %v\n", err)
			}
			return
		}

		message = strings.TrimSpace(message)
		fmt.Printf("收到来自 %s 的消息: %s\n", conn.RemoteAddr().String(), message)

		// 立即回复 "hello-2"
		response := "hello-2\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("发送回复失败: %v\n", err)
			return
		}
		fmt.Printf("回复 %s: %s", conn.RemoteAddr().String(), strings.TrimSpace(response))
	}
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		conn.Close()
	}
	s.conns = make(map[net.Conn]bool)

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
