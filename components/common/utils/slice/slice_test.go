// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slice provides a slice sorting function.
//
// It uses gross, low-level operations to make it easy to sort
// arbitrary slices with only a less function, without defining a new
// type with Len and Swap operations.
package utils

import (
	"fmt"
	"testing"
)

type User struct {
	Id  int64
	Age int
}

func Test_sorts(t *testing.T) {
	users := []*User{
		&User{Id: 1, Age: 10},
		&User{Id: 3, Age: 30},
		&User{Id: 2, Age: 20},
	}

	// 按Id升序
	Sort(users, func(i int, j int) bool {
		return users[i].Id < users[j].Id
	})
	for _, v := range users {
		fmt.Println("Id升序", v)
	}

	// 按Age降序
	Sort(users, func(i int, j int) bool {
		return users[i].Age > users[j].Age
	})
	for _, v := range users {
		fmt.Println("Age降序", v)
	}
}
