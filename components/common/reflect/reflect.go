package reflect

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func ReflectTry(f reflect.Value, args []reflect.Value, handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("-------------panic recover---------------")
			if handler != nil {
				handler(err)
			}
		}
	}()
	f.Call(args)
}

func GetStructName(v interface{}) string {
	return reflect.Indirect(reflect.ValueOf(v)).Type().Name()
}

func GetFuncName(fn interface{}) string {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("[fn = %v] is not func type.", fn))
	}

	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return cutLastString(fullName, ".", "-")
}

// cutLastString 截取字符串中最后一段，以@beginChar开始,@endChar结束的字符
// @text 文本
// @beginChar 开始
func cutLastString(text, beginChar, endChar string) string {
	if text == "" || beginChar == "" || endChar == "" {
		return ""
	}

	textRune := []rune(text)

	beginIndex := strings.LastIndex(text, beginChar)
	endIndex := strings.LastIndex(text, endChar)
	if endIndex < 0 || endIndex < beginIndex {
		endIndex = len(textRune)
	}

	return string(textRune[beginIndex+1 : endIndex])
}

func IsPtr(val interface{}) bool {
	if val == nil {
		return false
	}

	return reflect.TypeOf(val).Kind() == reflect.Ptr
}

func IsNotPtr(val interface{}) bool {
	if val == nil {
		return false
	}

	return reflect.TypeOf(val).Kind() != reflect.Ptr
}
