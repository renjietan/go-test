package pool

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	// GlobalPool 全局协程池实例
	GlobalPool *ants.Pool
	once       sync.Once
)

// InitPool 初始化协程池
// size: 协程池大小，0 表示使用默认值（默认是 math.MaxInt32）
func InitPool(size int) error {
	var err error
	once.Do(func() {
		if size <= 0 {
			// 使用默认大小
			GlobalPool, err = ants.NewPool(ants.DefaultAntsPoolSize)
		} else {
			GlobalPool, err = ants.NewPool(size)
		}
	})
	return err
}

// Submit 提交任务到协程池
func Submit(task func()) error {
	if GlobalPool == nil {
		return ants.ErrPoolOverload
	}
	return GlobalPool.Submit(task)
}

// Running 返回当前正在运行的协程数量
func Running() int {
	if GlobalPool == nil {
		return 0
	}
	return GlobalPool.Running()
}

// Free 返回当前空闲的协程数量
func Free() int {
	if GlobalPool == nil {
		return 0
	}
	return GlobalPool.Free()
}

// Cap 返回协程池的容量
func Cap() int {
	if GlobalPool == nil {
		return 0
	}
	return GlobalPool.Cap()
}

// Release 释放协程池
func Release() {
	if GlobalPool != nil {
		GlobalPool.Release()
		GlobalPool = nil
	}
}

// Reboot 重启协程池
func Reboot(size int) error {
	Release()
	return InitPool(size)
}
