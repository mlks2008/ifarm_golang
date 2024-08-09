/**
 * https_normal.go
 * ============================================================================
 * 一般使用方式
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

func normalOutput(w http.ResponseWriter, ret int64, desc string, data interface{}) error {

	var result = make(map[string]interface{})

	result["code"] = ret
	result["message"] = desc
	result["ts"] = time.Now().UnixNano() / 1e6
	result["data"] = data

	json, _ := json.Marshal(result)
	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

func normalInput(w http.ResponseWriter, r *http.Request) {

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
