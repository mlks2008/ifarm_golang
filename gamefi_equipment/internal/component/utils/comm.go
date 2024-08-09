package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func JoinInt64(arr []int64) string {
	var join string
	for _, a := range arr {
		join += fmt.Sprintf("%v,", a)
	}
	return strings.TrimRight(join, ",")
}

func JoinStr(arr []string) string {
	var join string
	for _, a := range arr {
		join += fmt.Sprintf("%v,", a)
	}
	return strings.TrimRight(join, ",")
}

func SplitToInt64(val string) []int64 {
	var ints = make([]int64, 0)

	var vals = strings.Split(val, ",")
	for _, v := range vals {
		i, _ := strconv.ParseInt(v, 10, 64)
		ints = append(ints, i)
	}
	return ints
}

func SplitToStr(val string) []string {
	return strings.Split(val, ",")
}
