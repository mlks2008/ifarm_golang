package strs

import (
	"fmt"
	"strings"
)

func Contains(s, substr string) bool {
	s = fmt.Sprintf(",%v,", s)
	substr = fmt.Sprintf(",%v,", substr)
	return strings.Contains(s, substr)
}
