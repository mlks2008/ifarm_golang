package go_demos

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

// 交替打印数字和字母
func Test_nlPrint(t *testing.T) {
	var number = make(chan struct{})
	var letter = make(chan struct{})
	var exit = make(chan struct{})

	go func() {

		var i = 1
		for {
			select {
			case _, ok := <-exit:
				if ok == false {
					return
				}
			case <-number:
				fmt.Print(i)
				i++
				fmt.Print(i)
				i++
				letter <- struct{}{}
			}
		}
	}()

	go func() {
		var i = 'A'
		for {
			select {
			case _, ok := <-exit:
				if ok == false {
					return
				}
			case <-letter:
				if i >= 'Z' {
					close(exit)
					continue
				}
				fmt.Print(string(i))
				i++
				fmt.Print(string(i))
				i++
				number <- struct{}{}
			}
		}
	}()

	number <- struct{}{}
	<-exit
}

// 给定一个string s1和一个string s2，请返回一个bool，代表两串是否重新排列后可相同。
func Test_reGroupString(t *testing.T) {
	var reGroupString = func(s1, s2 string) bool {
		var l1 = len([]rune(s1))
		var l2 = len([]rune(s2))
		if l1 > 5000 || l2 > 5000 || l1 != l2 {
			return false
		}
		var count = make(map[string]int)
		for _, v := range s1 {
			count[string(v)]++
		}
		for _, v := range s2 {
			count[string(v)]--
			if count[string(v)] < 0 {
				return false
			}
		}
		return true
	}

	fmt.Println(reGroupString("aa", "ab"))
}

// 翻转字符串
func Test_reverString(t *testing.T) {
	var reverString = func(s string) (string, error) {
		str := []rune(s)
		l := len(str)
		if l > 5000 {
			return s, fmt.Errorf("len out of range")
		}
		for i := 0; i < l/2; i++ {
			str[i], str[l-1-i] = str[l-1-i], str[i]
		}
		return string(str), nil
	}

	fmt.Println(reverString("abc"))
}

// 请编写一个方法，将字符串中的空格全部替换为“%20”。
// 假定该字符串有足够的空间存放新增的字符，并且知道字符串的真实长度(小于等于1000)，同时保证字符串由【大小写的英文字母组成】。
func Test_replaceBlank(t *testing.T) {
	var replaceBlank = func(s string) (string, error) {
		str := []rune(s)
		if len(str) > 1000 {
			return s, fmt.Errorf("len out of range")
		}
		for _, v := range str {
			if string(v) != " " && unicode.IsLetter(v) == false {
				return s, fmt.Errorf("nonletter")
			}
		}
		return strings.Replace(s, " ", "%20", -1), nil
	}

	fmt.Println(replaceBlank("abc d"))
}

const (
	Left = iota
	Top
	Right
	Bottom
)

func Test_bot(t *testing.T) {
	println(move("R2(LF)", 0, 0, Top))
}
func move(cmd string, x0 int, y0 int, z0 int) (x, y, z int) {
	x, y, z = x0, y0, z0
	repeat := 0
	repeatCmd := ""
	for _, s := range cmd {
		switch {
		case unicode.IsNumber(s):
			repeat = repeat*10 + (int(s) - '0')
		case s == ')':
			for i := 0; i < repeat; i++ {
				x, y, z = move(repeatCmd, x, y, z)
			}
			repeat = 0
			repeatCmd = ""
		case repeat > 0 && s != '(' && s != ')':
			repeatCmd = repeatCmd + string(s)
		case s == 'L':
			z = (z + 1) % 4
		case s == 'R':
			z = (z - 1 + 4) % 4
		case s == 'F':
			switch {
			case z == Left || z == Right:
				x = x - z + 1
			case z == Top || z == Bottom:
				y = y - z + 2
			}
		case s == 'B':
			switch {
			case z == Left || z == Right:
				x = x + z - 1
			case z == Top || z == Bottom:
				y = y + z - 2
			}
		}
	}
	return
}

func Test_code(t *testing.T) {

}
