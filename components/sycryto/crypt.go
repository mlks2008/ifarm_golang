package sycryto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"math/rand"
	"time"
)

func RandomAesKey(bitSize int) string {
	switch bitSize {
	case 16:
	case 24:
	case 32:
	default:
		bitSize = 16
	}
	return RandomKey(bitSize)
}

func RandomKey(bitSize int) string {
	keyBytes := make([]byte, 16)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 16; i++ {
		keyBytes[i] = byte(r.Intn(255))
	}
	return Bytes2String(keyBytes)
}

func Encrypt(text string, cipherkey string) (string, error) {
	key, err := String2Bytes(cipherkey)
	if err != nil {
		return "", err
	}

	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}

func Decrypt(encrypted string, cipherkey string) (string, error) {
	var err error
	key, err := String2Bytes(cipherkey)
	if err != nil {
		return "", err
	}

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}

func KeyEncryp(text string) (string, error) {
	return Encrypt(text, CIPHER_KEY)
}
func KeyDecrypt(encrypted string) (string, error) {
	return Decrypt(encrypted, CIPHER_KEY)
}

func Bytes2String(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func String2Bytes(str string) ([]byte, error) {
	return hex.DecodeString(str)
}
