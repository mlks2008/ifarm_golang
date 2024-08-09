package random

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

func Intn(val int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(val)
}

func InRange(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if max == min {
		max += 1
	}
	return rand.Intn(max-min) + min
}

func InFloatRange[F float32 | float64](min, max F) F {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	precisionNum := float64(10)
	minInt := int(min * F(precisionNum))
	maxInt := int(max*F(precisionNum)) + 1
	resInt := InRange(minInt, maxInt)

	return F(decimal.NewFromFloat(float64(resInt)).DivRound(decimal.NewFromInt(int64(precisionNum)), 2).InexactFloat64())
}

func WithNoRepeated(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn(end-start) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
