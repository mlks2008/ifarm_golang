package go_demos

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_search(t *testing.T) {
	var data = make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}
	var dataLen = len(data)
	var target = 438
	var size = 100

	var searchTarget = func(data []int, target int, ctx context.Context, result chan struct{}) {
		for _, v := range data {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if v == target {
				result <- struct{}{}
				return
			}
		}
	}

	var ctx, cancel = context.WithCancel(context.Background())
	var result = make(chan struct{})
	for i := 0; i < dataLen; i += size {
		end := i + size
		if end >= dataLen {
			end = dataLen - 1
		}
		go searchTarget(data[i:end], target, ctx, result)
	}

	var timer = time.NewTimer(time.Second * 5)
	select {
	case <-timer.C:
		cancel()
		fmt.Println("time out")
	case <-result:
		cancel()
		fmt.Println("find it")
	}
	time.Sleep(time.Second * 5)
}
