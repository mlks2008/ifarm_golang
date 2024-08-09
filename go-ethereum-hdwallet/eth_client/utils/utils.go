package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func IntToHex(num int64) string {
	value := uint64(num)
	return hexutil.EncodeUint64(value)
}

func Hex2Ten(hex string) (int64, error) {
	value, err := hexutil.DecodeUint64(hex)
	return int64(value), err
}

func Hex2Ten2(hex string) (string, error) {
	value, err := hexutil.DecodeBig(hex)
	return fmt.Sprint(value), err
}
