package go_demos

import (
	"fmt"
	"testing"
)

func add(x, y int) *int {
	res := x + y
	return &res
}

func Test_runaway(t *testing.T) {
	fmt.Println(add(1, 2))
}
