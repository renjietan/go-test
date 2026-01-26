package main

import (
	"fmt"
	"unsafe"
)

type Ss struct {
	Name string
}

func main() {

	a := make([]interface{}, 2)
	// a = append(a, 1)
	// a = append(a, "2")
	// a = append(a, map[string]string{
	// 	"3": "3",
	// })
	fmt.Println("容量:", cap(a))
	a[0] = 1
	a[1] = 2
	a = append(a, "3")
	a = append(a, map[string]string{
		"4": "4",
	})
	a = append(a, map[string]string{
		"5": "6",
	})
	a = append(a, map[string]string{
		"6": "6",
	})
	fmt.Println("容量:", cap(a))
	for index, v := range a {
		cc := uintptr(unsafe.Pointer(&v))
		fmt.Printf("%d v1: %v, 指针: %p  指针2: %#016x  指针3: %v \n", index, v, &v, &v, cc)
	}
	var b = make([]int, 3)
	b[0] = 1
	b[1] = 2
	b[2] = 3
	for index, v := range b {
		cc := uintptr(unsafe.Pointer(&v))
		fmt.Printf("%d b1: %v, 指针: %p  指针2: %#016x  指针3: %v \n", index, v, &v, &v, cc)
	}
	var aa = 1
	var bb = "2"
	var cc = "4"
	fmt.Printf("v1: %v, 指针: %p  指针2: %#016x  指针3: %v\n", aa, &aa, &aa, uintptr(unsafe.Pointer(&aa)))
	fmt.Printf("v1: %v, 指针: %p  指针2: %#016x  指针3: %v\n", bb, &bb, &bb, uintptr(unsafe.Pointer(&bb))-uintptr(unsafe.Pointer(&aa)))
	fmt.Printf("v1: %v, 指针: %p  指针2: %#016x  指针3: %v\n", cc, &cc, &cc, uintptr(unsafe.Pointer(&cc))-uintptr(unsafe.Pointer(&bb)))
}
