/**
 * http_command.go
 * ============================================================================
 * 执行请求
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:17
 */

package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

var G_HttpCommandM = NewHttpCommandM()

type HttpCommand interface {
	Prepare(http.ResponseWriter, *http.Request) bool
	Execute(json.RawMessage, http.ResponseWriter, *http.Request) bool
	Finish(http.ResponseWriter, *http.Request)
}

type HttpCommandM struct {
	cmdm map[string]HttpCommand
}

func NewHttpCommandM() *HttpCommandM {
	return &HttpCommandM{
		cmdm: make(map[string]HttpCommand),
	}
}

func (this *HttpCommandM) Register(name string, cmd HttpCommand) {
	this.cmdm[name] = cmd
}

func (this *HttpCommandM) Dispatcher(method string, params json.RawMessage, w http.ResponseWriter, r *http.Request) (bool, error) {
	if cmd, exist := this.cmdm[method]; exist {
		//每次新实例
		cmdType := reflect.TypeOf(cmd).Elem()
		cmdNew := reflect.New(cmdType).Interface().(HttpCommand)

		if cmdNew.Prepare(w, r) {
			ret := cmdNew.Execute(params, w, r)
			cmdNew.Finish(w, r)
			return ret, nil
		} else {
			return false, nil
		}
	}

	return false, errors.New(fmt.Sprintf("not found %s", method))
}
