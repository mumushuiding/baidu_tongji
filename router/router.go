package router

import (
	"net/http"

	"github.com/mumushuiding/baidu_tongji/config"
	"github.com/mumushuiding/baidu_tongji/controller"
	"github.com/mumushuiding/baidu_tongji/model"
)

// RouteFunction 根据路径指向方法
type RouteFunction func(*model.EditorTongji) (string, error)

// RouterMap 路由
var RouterMap map[string]RouteFunction

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
	Mux.HandleFunc("/check/memory", controller.MemoryCheck)
	Mux.HandleFunc("/api/v1/test/index", interceptor(controller.Index))
	// 获取统计数据接口
	Mux.HandleFunc("/api/v1/tongji/getData", interceptor(controller.GetTongjiData))
	Mux.HandleFunc("/api/v1/tongji/exportData", interceptor(controller.ExportData))

	// 百度统计接口
	Mux.HandleFunc("/api/v1/baidutongji/getData", interceptor(controller.GetBaiduDataByTimeSpan))
	// 远程拉取最新稿件
	Mux.HandleFunc("/api/v1/tongji/getFzmanuscript", interceptor(controller.GetRemoteFzManuscriptNotHave))
}
