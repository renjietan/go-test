package test_case

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRecover(t *testing.T) {
	var wg sync.WaitGroup
	var cn = make(chan int, 10)
	var map1 sync.Map
	map1 = sync.Map{}
	var map2 = make(map[int]int, 10)
	for v := range 10 {
		wg.Add(1)
		go func(index int) {
			defer func() {
				wg.Done()
				err := recover()
				if err != nil {
					fmt.Println("err", err)
				}
			}()
			cn <- index
			map1.Store(v, v)
			map2[v] = v
		}(v)
	}
	fmt.Println("当前时间1", time.Now().Format("2006-01-02 15:03:04"))
	wg.Wait()
	panic("err========================")
}
