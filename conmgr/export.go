package conmgr

import (
	"fmt"
	"reflect"

	"github.com/mumushuiding/baidu_tongji/model"
)

// ExportEditorFlowAndManuscript 导出指定时间段的流量和稿件
func ExportEditorFlowAndManuscript(e *model.EditorTongji) error {
	err := GetFlowAndManuscriptNum(e)
	if err != nil {
		return err
	}
	transform2Csv(e)
	return err
}

// ExportEditorFlowAndManuscriptNumLastMonth 导出上月编辑的流量
func ExportEditorFlowAndManuscriptNumLastMonth(e *model.EditorTongji) error {
	err := GetFlowAndManuscriptNumLastMonth(e)
	if err != nil {
		return err
	}
	transform2Csv(e)
	return err
}
func transform2Csv(e *model.EditorTongji) {
	// 显示的字段
	e.Body.Fields = []string{"用户名", "点击量", "稿件量", "访客数"}
	// e.Body.Fields = []string{"点击量", "稿件量", "访客数"}
	// 获取结果集
	s := reflect.ValueOf(e.Body.Data[0])
	// 遍历结果
	for i := 0; i < s.Len(); i++ {
		ele := s.Index(i)
		flow := ele.Interface().(*model.EURLFlow)
		var list []string
		list = append(list, flow.Realname)
		list = append(list, fmt.Sprintf("%d", flow.PvCount))
		list = append(list, fmt.Sprintf("%d", flow.Manuscript))
		list = append(list, fmt.Sprintf("%d", flow.VisitorCount))
		e.Body.Data = append(e.Body.Data, list)
	}
	e.Body.Data = e.Body.Data[1:]
}
