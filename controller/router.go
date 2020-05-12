package controller

import (
	"errors"

	"github.com/mumushuiding/baidu_tongji/conmgr"
	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/baidu_tongji/service"
)

// RouteFunction 根据路径指向方法
type RouteFunction func(*model.EditorTongji) (string, error)

// RouterMap 路由
var RouterMap map[string]RouteFunction

// SetRouterMap 设置函数路由
func SetRouterMap() {
	RouterMap = make(map[string]RouteFunction)
	RouterMap["visit/editor/flow"] = service.FindAllEditorFlowPaged
	RouterMap["visit/editor/flowWithAvators"] = service.FindAllEditorFlowPagedWithAvators
	RouterMap["visit/editor/details"] = service.FindEditorDetails
	RouterMap["visit/editor/articles"] = service.FindEditorArticles
	RouterMap["visit/editor/trend/visitor"] = service.FindEditorTrendVisitor
	RouterMap["visit/article/flowWithAvators"] = conmgr.GetArticleFlowWithAvators
	RouterMap["visit/editor/flowAndManuscriptNumLastMonth"] = conmgr.GetFlowAndManuscriptNumLastMonth
}

// GetRoute 获取执行函数
func GetRoute(route string) (func(*model.EditorTongji) (string, error), error) {
	f := RouterMap[route]
	if f == nil {
		return nil, errors.New("method:" + route + ",不存在")
	}
	return f, nil
}
