/**
 * https_httpapi.go
 * ============================================================================
 * apikey对外接接口格式
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:17
 */

package httpserver

import (
	"encoding/json"
	"fmt"
	"goarbitrage/pkg/httpserver/retno"
	"net/http"
	"strings"
	"time"
)

type HttpApiJsonResult struct {
	Id     int64         `json:"id,omitempty"`
	Ts     int64         `json:"ts"`
	Error  *HttpApiError `json:"error,omitempty"`
	Result interface{}   `json:"result"`
}

type HttpApiError struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

func httpApiOutput(w http.ResponseWriter, id int64, ret int64, desc string, data interface{}) error {

	var result = HttpApiJsonResult{
		Id:     id,
		Ts:     time.Now().UnixNano() / 1e6,
		Result: data,
	}

	if ret != retno.SYS_OK {
		result.Error = &HttpApiError{Code: ret, Msg: desc}
	}

	json, _ := json.Marshal(result)

	if _, err := w.Write(json); err != nil {
		//l4g.Error(err)
		return err
	}

	return nil
}

func httpApiInput(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var ret bool

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

	if false {
		if ret {
			//l4g.Debug("ApiHandler info is %+v done", path)
		} else {
			//l4g.Error("ApiHandler info is %+v fail", path)
		}
	}
}
