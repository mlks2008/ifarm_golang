package utils

func PrintPanicStack(extras ...interface{}) bool {
	if x := recover(); x != nil {
		//l4g.Error(x)
		//i := 0
		//funcName, file, line, ok := runtime.Caller(i)
		//for ok {
		//	l4g.Error("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
		//	i++
		//	funcName, file, line, ok = runtime.Caller(i)
		//}
		//
		//for k := range extras {
		//	l4g.Error("EXRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		//}

		return true
	} else {
		return false
	}
}
