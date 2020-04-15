package model

import (
	"time"

	"github.com/mumushuiding/util"
)

// EditorTongji 用户统计数据
type EditorTongji struct {
	Header BHeader `json:"header"`
	Body   EBody   `json:"body"`
}

// EBody EBody
type EBody struct {
	Data       []interface{} `json:"data,omitempty"`
	Total      int           `json:"total,omitempty"`
	StartIndex int           `json:"start_index,omitempty"`
	MaxResults int           `json:"max_results,omitempty"`
	StartDate  string        `json:"start_date,omitempty"`
	EndDate    string        `json:"end_date,omitempty"`
	Method     string        `json:"method"`
}

// EURLFlow 编辑对应的流量
type EURLFlow struct {
	// username 唯一
	Username string `json:"username"`
	// realname不唯一
	Realname        string `json:"realname"`
	PvCount         int    `json:"pv_count,omitempty"`
	VisitorCount    int    `json:"visitor_count,omitempty"`
	OutwardCount    int    `json:"outward_count,omitempty"`
	ExitCount       int    `json:"exit_count,omitempty"`
	AverageStayTime int    `json:"average_stay_time,omitempty"`
}

// ToString ToString
func (e *EditorTongji) ToString() string {
	str, _ := util.ToJSONStr(e)
	return str
}

// FindAllEditorFlowPaged 分页查询所有的编辑流量
// select e.username,e.realname,sum(f.pv_count) pv_count from baidu_url_editor e join baidu_url_flow f on e.page_id=f.page_id and time_span>='2020-04-05' and time_span<='2020-04-07
// ' group by e.username,e.realname order by pv_count desc  limit 30;
func (e *EditorTongji) FindAllEditorFlowPaged() (string, error) {
	euf, err := e.findAllEditor()
	if err != nil {
		return "", err
	}
	e.Body.Data = append(e.Body.Data, euf)
	return e.ToString(), nil
}

func (e *EditorTongji) findAllEditor() ([]*EURLFlow, error) {
	var result []*EURLFlow
	if len(e.Body.StartDate) == 0 {
		e.Body.StartDate = util.FormatDate3(time.Now())
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = e.Body.StartDate
	}
	if e.Body.MaxResults == 0 {
		e.Body.MaxResults = 20
	}
	var total int
	err := db.Table("baidu_url_editor").
		Select("baidu_url_editor.username,baidu_url_editor.realname,sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time").
		Joins("left join baidu_url_flow on baidu_url_flow.page_id=baidu_url_editor.page_id").
		Group("baidu_url_editor.username,baidu_url_editor.realname").
		Where("baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", e.Body.StartDate, e.Body.EndDate).
		Count(&total).
		Order("pv_count desc,visitor_count desc").
		Offset(e.Body.StartIndex).
		Limit(e.Body.MaxResults).Find(&result).Error
	e.Body.Total = total
	return result, err
}
