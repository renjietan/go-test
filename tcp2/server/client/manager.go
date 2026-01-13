package client

import (
	"net"
	"sync"
	"sync/atomic"
)

// Info 客户端信息
type Info struct {
	ID       int64    // 客户端唯一ID
	Conn     net.Conn // 连接对象
	Addr     string   // 客户端地址
	IsActive bool     // 是否活跃
}

// Manager 客户端管理器
type Manager struct {
	clients sync.Map // map[int64]*Info
	counter int64    // 客户端ID计数器
}

var (
	// GlobalManager 全局客户端管理器实例，可在程序任意一处访问
	GlobalManager = &Manager{}
)

// AddClient 添加客户端
func (cm *Manager) AddClient(conn net.Conn) *Info {
	id := atomic.AddInt64(&cm.counter, 1)
	client := &Info{
		ID:       id,
		Conn:     conn,
		Addr:     conn.RemoteAddr().String(),
		IsActive: true,
	}
	cm.clients.Store(id, client)
	return client
}

// RemoveClient 移除客户端
func (cm *Manager) RemoveClient(id int64) {
	cm.clients.Delete(id)
}

// GetClient 获取客户端
func (cm *Manager) GetClient(id int64) (*Info, bool) {
	value, ok := cm.clients.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*Info), true
}

// GetAllClients 获取所有客户端
func (cm *Manager) GetAllClients() []*Info {
	var clients []*Info
	cm.clients.Range(func(key, value interface{}) bool {
		client := value.(*Info)
		if client.IsActive {
			clients = append(clients, client)
		}
		return true
	})
	return clients
}

// GetClientCount 获取客户端数量
func (cm *Manager) GetClientCount() int {
	count := 0
	cm.clients.Range(func(key, value interface{}) bool {
		client := value.(*Info)
		if client.IsActive {
			count++
		}
		return true
	})
	return count
}

// Broadcast 广播消息给所有客户端
func (cm *Manager) Broadcast(message string) {
	cm.clients.Range(func(key, value interface{}) bool {
		client := value.(*Info)
		if client.IsActive {
			_, err := client.Conn.Write([]byte(message))
			if err != nil {
				client.IsActive = false
				cm.RemoveClient(client.ID)
			}
		}
		return true
	})
}
