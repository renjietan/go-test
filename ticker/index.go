package main

import (
	"fmt"
	"time"
)

func main() {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("超时")
	default:
	}
	time.Sleep(time.Second * 10)
}
