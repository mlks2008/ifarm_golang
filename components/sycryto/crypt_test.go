package sycryto

import (
	"fmt"
	"testing"
)

func Test_EncryptDecrypt(t *testing.T) {
	var cipher = RandomAesKey(16)
	enDate, err := Encrypt("test", cipher)
	fmt.Println(enDate, err)
	deDate, err := Decrypt(enDate, cipher)
	fmt.Println(deDate, err)
}
