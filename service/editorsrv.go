package service

import "github.com/mumushuiding/baidu_tongji/model"

// FindAllEditorFlowPaged 查询编辑的流量
func FindAllEditorFlowPaged(e *model.EditorTongji) (string, error) {
	return e.FindAllEditorFlowPaged()
}
