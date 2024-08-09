package pprof

import (
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
	"net/http/pprof"
)

func RunPprofServer(addr string) {
	pprofServer := http.NewServeMux()
	pprofServer.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pprofServer.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pprofServer.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pprofServer.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pprofServer.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	go func() {
		if e := http.ListenAndServe(addr, pprofServer); e != nil {
			log.Errorf("pprof server error ListenAndServe addr:%s,error:%+v", addr, e)
		} else {
			log.Infof("pprof server ListenAndServe addr:%s success", addr)
		}
		defer func() {
			if e := recover(); e != nil {
				log.Errorf("expected panic from pprof server,error:%+v", e)
			}
		}()
	}()
}
