/**
 * http_handler_pool.go
 * ============================================================================
 * 请求池
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:17
 */

package httpserver

import (
	"encoding/json"
	. "goarbitrage/pkg/httpserver/register"
	"goarbitrage/pkg/httpserver/retno"
	"net/http"
)

type HttpRequestInfo struct {
	method string
	params json.RawMessage
	w      http.ResponseWriter
	r      *http.Request
	close  chan bool
}

func NewHttpRequestInfo(method string, params json.RawMessage, w http.ResponseWriter, r *http.Request, c chan bool) *HttpRequestInfo {
	req := &HttpRequestInfo{
		method: method,
		params: params,
		w:      w,
		r:      r,
		close:  c,
	}
	return req
}

type HttpHandlerPool struct {
	pool_id int
}

func NewHttpHandlerPool(i int) *HttpHandlerPool {
	return &HttpHandlerPool{
		pool_id: i,
	}
}
func (this *HttpHandlerPool) Process() {
	for {
		//l4g.Debug("HttpHandlerPool processs %d......", this.pool_id)
		select {
		case request := <-g_handler_chan:

			//l4g.Debug("HttpHandlerPool info is %+v, processs is %d, content-type: %v, Form: %+v PostForm: %+v", request.method, this.pool_id, request.r.Header.Get("Content-Type"), request.r.Form, request.r.PostForm)

			ret, err := G_HttpCommandM.Dispatcher(request.method, request.params, request.w, request.r)
			if err != nil {
				HttpOutput(request.w, retno.SYS_METHOD_NOT_FIND, err.Error(), nil)
			}

			if ret {
				//l4g.Debug("HttpHandlerPool info is %+v, processs is %d, success", request.method, this.pool_id)
			} else {
				//l4g.Error("HttpHandlerPool info is %+v, processs is %d, fail", request.method, this.pool_id)
			}

			request.close <- ret
		}
	}
	//l4g.Debug("HttpHandlerPool processs %d finish", this.pool_id)
}
