package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_code(t *testing.T) {
	now := time.Now().Unix()
	old := now - 24*3600 + 1
	fmt.Println((now - old) / (24 * 3600))
}
