package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.RWMutex
var count int

func main() {
	go A()
	time.Sleep(2 * time.Second)
	go func() {
		fmt.Println("lock1")
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("lock2")
		count++
		fmt.Println(count)
	}()
	time.Sleep(time.Second * 10)
}
func A() {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("a")
	B()
}
func B() {
	fmt.Println("b")
	time.Sleep(5 * time.Second)
	C()
}
func C() {
	fmt.Println("c")
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("cc")
}
