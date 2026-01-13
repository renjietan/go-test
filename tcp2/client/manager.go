package client

import (
	"sync"
	"sync/atomic"
)

// ClientInfo TCP客户端信息
type ClientInfo struct {
	ID       int64  // 客户端唯一ID
	Addr     string // 服务端地址
	IsActive bool   // 是否活跃
	IsConnected bool // 是否已连接
}

// Manager 客户端管理器
type Manager struct {
	clients sync.Map // map[int64]*ClientInfo
	counter int64    // 客户端ID计数器
	mu      sync.RWMutex
}

var (
	// GlobalManager 全局客户端管理器实例，可在程序任意一处访问
	GlobalManager = &Manager{}
	
	// GlobalClients 全局客户端实例存储，可在程序任意一处访问
	GlobalClients sync.Map // map[int64]*TCPClient
)

// AddClient 添加客户端
func (cm *Manager) AddClient(addr string) *ClientInfo {
	id := atomic.AddInt64(&cm.counter, 1)
	client := &ClientInfo{
		ID:          id,
		Addr:        addr,
		IsActive:    true,
		IsConnected: false,
	}
	cm.clients.Store(id, client)
	return client
}

// UpdateClientStatus 更新客户端状态
func (cm *Manager) UpdateClientStatus(id int64, isConnected bool) {
	if value, ok := cm.clients.Load(id); ok {
		client := value.(*ClientInfo)
		cm.mu.Lock()
		client.IsConnected = isConnected
		cm.mu.Unlock()
	}
}

// RemoveClient 移除客户端
func (cm *Manager) RemoveClient(id int64) {
	cm.clients.Delete(id)
}

// GetClient 获取客户端
func (cm *Manager) GetClient(id int64) (*ClientInfo, bool) {
	value, ok := cm.clients.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*ClientInfo), true
}

// GetAllClients 获取所有客户端
func (cm *Manager) GetAllClients() []*ClientInfo {
	var clients []*ClientInfo
	cm.clients.Range(func(key, value interface{}) bool {
		client := value.(*ClientInfo)
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
		client := value.(*ClientInfo)
		if client.IsActive {
			count++
		}
		return true
	})
	return count
}

// GetConnectedCount 获取已连接的客户端数量
func (cm *Manager) GetConnectedCount() int {
	count := 0
	cm.clients.Range(func(key, value interface{}) bool {
		client := value.(*ClientInfo)
		if client.IsActive && client.IsConnected {
			count++
		}
		return true
	})
	return count
}

// StoreClient 存储客户端实例（供外部调用）
func StoreClient(id int64, client *TCPClient) {
	GlobalClients.Store(id, client)
}

// GetClientInstance 获取客户端实例（供外部调用）
func GetClientInstance(id int64) (*TCPClient, bool) {
	value, ok := GlobalClients.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*TCPClient), true
}

// GetAllClientInstances 获取所有客户端实例（供外部调用）
func GetAllClientInstances() []*TCPClient {
	var result []*TCPClient
	GlobalClients.Range(func(key, value interface{}) bool {
		result = append(result, value.(*TCPClient))
		return true
	})
	return result
}

// RemoveClientInstance 移除客户端实例
func RemoveClientInstance(id int64) {
	GlobalClients.Delete(id)
}
