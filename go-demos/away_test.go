package go_demos

import (
	"fmt"
	"testing"
)

func Test_away(t *testing.T) {
	var add = func(x, y int) *int {
		res := x + y
		return &res
	}
	fmt.Println(add(1, 2))
}
