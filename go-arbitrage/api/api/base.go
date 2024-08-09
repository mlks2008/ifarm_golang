/**
 * v.go.go
 * ============================================================================
 * 简介
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/3/2 12:13
 */

package api

import (
	. "goarbitrage/api"
	"net/http"
)

type Base struct {
	Api
}

// 执行Execute前，会先执行Prepare
func (this *Base) Prepare(w http.ResponseWriter, r *http.Request) bool {

	return true
}

// 执行Execute后，执行Finish
func (this *Base) Finish(w http.ResponseWriter, r *http.Request) {

}
