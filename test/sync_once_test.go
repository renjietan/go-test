package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestSyncOnce(t *testing.T) {
	var so sync.Once

	so.Do(func() {
		test("str1")
	})
	test("str2")
	t.Log("SyncOnceTest")
}

func test(str string) {
	fmt.Println(str)
}
