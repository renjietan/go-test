package test_case

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestPtr(t *testing.T) {
	// 查看几个指针地址
	vars := []interface{}{1, "hello", 3.14, true}

	for i, v := range vars {
		// 每个变量都有不同的地址
		fmt.Printf("变量%d 类型: %T, 地址: %p 地址(补0): %#016x\n", i, v, &v, &v)
	}

	// 内存对齐的影响
	var a, b, c int
	fmt.Printf("\n连续变量的地址:\n")
	fmt.Printf("a: %p\n", &a)
	fmt.Printf("b: %p\n", &b)
	fmt.Printf("c: %p\n", &c)

	addrA := uintptr(unsafe.Pointer(&a))
	addrB := uintptr(unsafe.Pointer(&b))
	fmt.Println("&a:", unsafe.Pointer(&a))
	fmt.Println("&b:", unsafe.Pointer(&b))
	fmt.Printf("a和b的地址差: %d 字节\n", addrB-addrA)
}
