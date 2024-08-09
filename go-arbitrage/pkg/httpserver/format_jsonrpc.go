/**
 * https_jsonrpc.go
 * ============================================================================
 * jsonrpc请求格式
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:17
 */

package httpserver

import (
	"encoding/json"
	"fmt"
	"goarbitrage/pkg/httpserver/retno"
	"goarbitrage/pkg/httpserver/utils"
	"net/http"
	"strings"
)

type JsonBody struct {
	Id     int64           `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type JsonResult struct {
	Id     int64       `json:"id"`
	Error  *JsonError  `json:"error,omitempty"`
	Result interface{} `json:"result"`
}

type JsonError struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

func jsonRpcOutput(w http.ResponseWriter, id int64, ret int64, desc string, data interface{}) error {

	var result JsonResult

	if ret != retno.SYS_OK {
		result = JsonResult{
			Id:     id,
			Error:  &JsonError{Code: ret, Msg: desc},
			Result: data,
		}
	} else {
		result = JsonResult{
			Id:     id,
			Result: data,
		}
	}

	json, _ := json.Marshal(result)

	//l4g.Debug("ApiHandler respone %s", json)

	if _, err := w.Write(json); err != nil {
		//l4g.Error(err)
		return err
	}

	return nil
}

func jsonRpcInput(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var jsonreq JsonBody

	var ret bool

	if r.Header.Get("Content-Type") == "application/json" {

		err := utils.UnDeserialize(r.Body, &jsonreq)
		if err != nil {
			//l4g.Error("ApiHandler info is %+v err %v", jsonreq, err)
			HttpOutput(w, retno.SYS_ERROR, err.Error(), nil)
			return
		}

		r.PostForm.Set("sys_api_method", jsonreq.Method)

		//l4g.Debug("ApiHandler info is %+v ", jsonreq)

		closeChan := make(chan bool, 0)

		g_handler_chan <- NewHttpRequestInfo(jsonreq.Method, jsonreq.Params, w, r, closeChan)

		ret = <-closeChan

	} else {

		path := strings.Trim(r.URL.Path, "/")

		if len(strings.Split(path, "/")) < 2 {
			//l4g.Error("ApiHandler path is invalid:%v", r.URL.Path)
			HttpOutput(w, retno.SYS_PATH_INVALID, fmt.Sprintf("path is invalid:%v", r.URL.Path), nil)
			return
		}

		r.PostForm.Set("sys_api_method", path)

		//l4g.Debug("ApiHandler info is %s ", path)

		closeChan := make(chan bool, 0)

		g_handler_chan <- NewHttpRequestInfo(path, nil, w, r, closeChan)

		ret = <-closeChan
	}

	if false {
		if ret {
			//l4g.Debug("ApiHandler info is %+v done", jsonreq)
		} else {
			//l4g.Error("ApiHandler info is %+v fail", jsonreq)
		}
	}
}
