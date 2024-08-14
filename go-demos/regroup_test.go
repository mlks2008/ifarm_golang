package go_demos

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Regroup(t *testing.T) {
	var isRegroup = func(s1, s2 string) bool {
		sl1 := len([]rune(s1))
		sl2 := len([]rune(s2))

		if sl1 > 5000 || sl2 > 5000 || sl1 != sl2 {
			return false
		}

		for _, v := range s1 {
			if strings.Count(s1, string(v)) != strings.Count(s2, string(v)) {
				return false
			}
		}
		return true
	}

	fmt.Println(isRegroup("ab", "ba"))
}

func Test_rune(t *testing.T) {
	// 定义一个字符串，包含英文和中文字符
	s := "Hello, 世界"

	// 使用 for range 遍历字符串
	for i, r := range s {
		// r 的类型是 rune
		fmt.Printf("字符 %q 的索引是 %d，Unicode 码点是 %U\n", r, i, r)
	}

	// 将字符串转换为 []rune 切片
	runes := []rune(s)
	fmt.Println("字符串的 rune 切片是:", runes)

	// 通过索引访问 rune 切片
	fmt.Printf("第一个字符是 %q，Unicode 码点是 %U\n", runes[0], runes[0])
}
