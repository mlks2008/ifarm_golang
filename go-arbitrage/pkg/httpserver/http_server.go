/**
 * http_server.go
 * ============================================================================
 * http服务
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/7/4 15:35
 */

package httpserver

import (
	"fmt"
	"net/http"
)

var (
	g_httpformat   HttpFormat
	g_handler_chan = make(chan *HttpRequestInfo, 10240)
)

func RunApiServer(httpConns int, httpPerHostConns int, httpPoolNum int, listenAddr string) {

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = httpPerHostConns
	http.DefaultTransport.(*http.Transport).MaxIdleConns = httpConns

	//http pool
	for i := 0; i < httpPoolNum; i++ {
		pool := NewHttpHandlerPool(i)
		go pool.Process()
	}

	HttpRouter()

	fmt.Println(http.ListenAndServe(listenAddr, nil))
}

func SetHttpFormat(httpformat string) {
	g_httpformat = HttpFormat(httpformat)
}

func SetHandlerChan(size int) {
	g_handler_chan = make(chan *HttpRequestInfo, size)
}
