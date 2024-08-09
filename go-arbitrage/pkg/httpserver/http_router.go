/**
 * http_router.go
 * ============================================================================
 * 路由
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:17
 */

package httpserver

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpFormat string

var (
	JSONRPC HttpFormat = "jsonrpc" //标准jsonrpc格式请求响应
	NORMAL  HttpFormat = "normal"  //常规&格式get,post请求,json格式响应
	HTTPAPI HttpFormat = "httpapi" //公开API格式,参考火币格式
)

func HttpRouter() {
	http.HandleFunc("/", ApiHandler)
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	HttpInput(w, r)
}

//HTTP请求(默认为normal)
func HttpInput(w http.ResponseWriter, r *http.Request) {

	if g_httpformat == JSONRPC {
		jsonRpcInput(w, r)
		return
	}

	if g_httpformat == HTTPAPI {
		httpApiInput(w, r)
		return
	}

	if g_httpformat == NORMAL || g_httpformat == "" {
		normalInput(w, r)
		return
	}
}

//HTTP响应(默认为normal)
func HttpOutput(w http.ResponseWriter, ret int64, desc string, data interface{}) error {

	if g_httpformat == JSONRPC {
		return jsonRpcOutput(w, 0, ret, desc, data)
	}

	if g_httpformat == HTTPAPI {
		return httpApiOutput(w, 0, ret, desc, data)
	}

	if g_httpformat == NORMAL || g_httpformat == "" {
		return normalOutput(w, ret, desc, data)
	}

	return errors.New(fmt.Sprintf("HttpOutput httpformat no exist:%v", g_httpformat))
}
