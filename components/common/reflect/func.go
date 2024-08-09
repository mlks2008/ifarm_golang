package reflect

import (
	"errors"
	"reflect"
)

var (
	nilFuncInfo = FuncInfo{}
)

type FuncInfo struct {
	Type       reflect.Type
	Value      reflect.Value
	InArgs     []reflect.Type
	InArgsLen  int
	OutArgs    []reflect.Type
	OutArgsLen int
}

func GetFuncInfo(fn interface{}) (FuncInfo, error) {
	if fn == nil {
		return nilFuncInfo, errors.New("function is nil")
	}

	typ := reflect.TypeOf(fn)

	if typ.Kind() != reflect.Func {
		return nilFuncInfo, errors.New("is not func type")
	}

	var inArgs []reflect.Type
	for i := 0; i < typ.NumIn(); i++ {
		t := typ.In(i)
		inArgs = append(inArgs, t)
	}

	var outArgs []reflect.Type
	for i := 0; i < typ.NumOut(); i++ {
		t := typ.Out(i)
		outArgs = append(outArgs, t)
	}

	funcInfo := FuncInfo{
		Type:       typ,
		Value:      reflect.ValueOf(fn),
		InArgs:     inArgs,
		InArgsLen:  typ.NumIn(),
		OutArgs:    outArgs,
		OutArgsLen: typ.NumOut(),
	}

	return funcInfo, nil
}
