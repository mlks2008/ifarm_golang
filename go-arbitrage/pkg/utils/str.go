/**
 * str.go
 * ============================================================================
 * 简介
 * ============================================================================
 * author: peter.wang
 * createtime: 2020-07-07 23:08
 */

package utils

import "math/rand"

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
