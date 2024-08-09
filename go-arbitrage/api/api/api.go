package api

import (
	"encoding/json"
	"goarbitrage/api"
	. "goarbitrage/pkg/httpserver"
	"goarbitrage/pkg/httpserver/retno"
	"net/http"
)

type Api_StopOnePlat struct{ Base }

func (this *Api_StopOnePlat) Execute(params json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
	var val bool
	err := json.Unmarshal(params, &val)
	if err != nil {
		HttpOutput(w, retno.SYS_ERROR, "err", err)
		return true
	} else {
		api.OnePlatStop = val
		HttpOutput(w, retno.SYS_OK, "succ", api.GetPrintInfo())
		return true
	}
}

type Api_SetOnePlat struct{ Base }

func (this *Api_SetOnePlat) Execute(params json.RawMessage, w http.ResponseWriter, r *http.Request) bool {
	type Param struct {
		OnePlatUsdtAmount1 float64 `json:"OnePlatUsdtAmount1"`
		OnePlatUsdtAmount2 float64 `json:"OnePlatUsdtAmount2"`
		OnePlatUsdtAmount3 float64 `json:"OnePlatUsdtAmount3"`
	}
	var param Param
	err := json.Unmarshal(params, &param)
	if err != nil {
		HttpOutput(w, retno.SYS_ERROR, "err", err)
		return true
	} else {
		api.OnePlatUsdtAmount1 = param.OnePlatUsdtAmount1
		api.OnePlatUsdtAmount2 = param.OnePlatUsdtAmount2
		api.OnePlatUsdtAmount3 = param.OnePlatUsdtAmount3
		HttpOutput(w, retno.SYS_OK, "succ", api.GetPrintInfo())
		return true
	}
}
