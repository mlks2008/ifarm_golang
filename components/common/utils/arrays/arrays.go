package arrays

import (
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// 从切片中根据索引位置，移除一些元素
func RemoveInSlice[T any](removeEleIdx []int, list []T) []T {
	pos := 0
	for i, _ := range list {
		if !lo.Contains(removeEleIdx, i) {
			list[pos] = list[i]
			pos++
		}
	}
	return list[:pos]
}

func SliceConvert2String(s1 []any) []string {
	s2 := make([]string, len(s1))
	for i := 0; i < len(s1); i++ {
		s2[i] = cast.ToString(s1[i])
	}
	return s2
}

type Number interface {
	~int | ~uint | ~int32 | ~uint32 | ~int64 | ~uint64 | ~int8 | ~uint8 | ~float64 | ~float32
}

func IntInterConvert[T1, T2 Number](s []T1) []T2 {
	s2 := make([]T2, len(s))
	for i := 0; i < len(s); i++ {
		s2[i] = T2(s[i])
	}
	return s2
}
