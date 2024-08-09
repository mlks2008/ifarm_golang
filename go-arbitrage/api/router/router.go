/**
 * router.go
 * ============================================================================
 * 简介
 * ============================================================================
 * author: peter.wang
 * createtime: 2018/11/19 14:20
 */

package httpapi

import (
	"goarbitrage/api/api"
	. "goarbitrage/pkg/httpserver/register"
)

func init() {
	G_HttpCommandM.Register("StopOnePlat", &api.Api_StopOnePlat{})
	G_HttpCommandM.Register("SetOnePlat", &api.Api_SetOnePlat{})
}
