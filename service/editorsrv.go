package service

import (
	"github.com/mumushuiding/baidu_tongji/model"
)

// FindAllEditorFlowPaged 查询编辑的流量
func FindAllEditorFlowPaged(e *model.EditorTongji) (string, error) {
	return e.FindAllEditorFlowPaged()
}

// FindAllEditorFlowPagedWithAvators FindAllEditorFlowPagedWithAvators
func FindAllEditorFlowPagedWithAvators(e *model.EditorTongji) (string, error) {
	err := e.FindAllEditorFlowPagedWithAvatar()
	if err != nil {
		return "", err
	}
	return e.ToString(), nil
}

// FindEditorArticles 查询编辑的文章
func FindEditorArticles(e *model.EditorTongji) (string, error) {
	err := e.FindEditorArticles()
	if err != nil {
		return "", err
	}
	return e.ToString(), nil
}

// FindEditorDetails 查询编辑所有信息
func FindEditorDetails(e *model.EditorTongji) (string, error) {
	err := e.FindEditorDetails()
	if err != nil {
		return "", err
	}
	return e.ToString(), nil
}

// FindEditorTrendVisitor 查询编辑对应的新老用户访问趋势
func FindEditorTrendVisitor(e *model.EditorTongji) (string, error) {
	err := e.FindEditorTrendVisitor()
	if err != nil {
		return "", err
	}
	return e.ToString(), nil
}
