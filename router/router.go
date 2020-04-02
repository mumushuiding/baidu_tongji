package router

import (
	"net/http"

	"github.com/mumushuiding/baidu_tongji/config"

	"github.com/mumushuiding/baidu_tongji/controller"
)

// Mux 路由
var Mux = http.NewServeMux()
var conf = *config.Config

func interceptor(h http.HandlerFunc) http.HandlerFunc {
	return crossOrigin(h)
}
func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", conf.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AccessControlAllowHeaders)
		h(w, r)
	}
}
func init() {
	setMux()
}
func setMux() {
	Mux.HandleFunc("/api/v1/test/index", interceptor(controller.Index))

	// 百度统计接口
	Mux.HandleFunc("/api/v1/baidutongji/getData", interceptor(controller.GetBaiduDataByTimeSpan))
}
