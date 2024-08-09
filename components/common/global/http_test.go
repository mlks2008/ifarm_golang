package global

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func Test_JwtToken(t *testing.T) {
	jwtToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1aWQiOiIxNTAyNDY1NDkyNjIzNDI5IiwiZXhwIjoyNjA1MzExMTQzfQ.x03dhxK2QMiKmOTkMb8dljD66ZEXkkeGxZlmJKyPdlw"
	tokens := strings.Split(jwtToken, ".")

	bb1, err := base64.RawStdEncoding.DecodeString(tokens[0])
	fmt.Println(fmt.Sprintf("%s", bb1), err)

	bb2, err := base64.RawStdEncoding.DecodeString(tokens[1])
	fmt.Println(fmt.Sprintf("%s", bb2), err)

	//bb3, err := base64.RawStdEncoding.DecodeString(tokens[2])
	//fmt.Println(fmt.Sprintf("%s", bb3), err)
}
