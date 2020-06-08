package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	UserName   string        `json:"username,omitempty"`
	Method     string        `json:"method,omitempty"`
	Metrics    string        `json:"metrics,omitempty"`
	Fields     []string      `json:"fields,omitempty"`
}

// EURLFlow 编辑对应的流量
// 平均访问页数
// 跳出率: ExitCount/PvCount 跳出率
type EURLFlow struct {
	// username 唯一
	Username string `json:"username,omitempty"`
	// realname不唯一
	Realname     string `json:"realname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	TimeSpan     string `json:"time_span,omitempty"`
	PageID       string `json:"page_id,omitempty"`
	Name         string `json:"name,omitempty"`
	PvCount      int    `json:"pv_count,omitempty"`
	VisitorCount int    `json:"visitor_count,omitempty"`
	// 新老用户比
	VisitorRatio string `json:"visitor_ratio,omitempty"`
	// OutwardCount 贡献下游流量
	OutwardCount int `json:"outward_count,omitempty"`
	// ExitCount 跳出次数 ExitCount/PvCount 为跳出率
	ExitCount int `json:"exit_count,omitempty"`
	// ExitCount/PvCount 为跳出率
	ExitRatio       string  `json:"exit_ratio,omitempty"`
	AverageStayTime int     `json:"average_stay_time,omitempty"`
	Title           string  `json:"title,omitempty"`
	Source          string  `json:"source,omitempty"`
	Star            float64 `json:"star,omitempty"`
	// 稿件数
	Manuscript int `json:"manuscript,omitempty"`
}

// ToString ToString
func (e *EditorTongji) ToString() string {
	str, _ := util.ToJSONStr(e)
	return str
}

// GetInfluence 获取传播力
func (u *EURLFlow) GetInfluence(startDate, endDate string) (float64, error) {
	// pv占比40%，visitor_count 60%,平均每日十万加为3星
	pvbase := 8000
	pv := 0.4
	visitor := 0.4
	OutwardCount := 0.2
	days, err := util.DateStrSubDays3(startDate, endDate)
	days++
	if err != nil {
		return 0, nil
	}
	// log.Println("s:", float64(u.PvCount)*pv+float64(u.VisitorCount)*visitor)
	// fmt.Printf("pv:%f,uv:%f,days:%d,startDate:%s,endDate:%s", float64(u.PvCount), float64(u.VisitorCount), days, startDate, endDate)
	result := 6 * (float64(u.PvCount)*pv + float64(u.VisitorCount)*visitor + float64(u.OutwardCount)*OutwardCount*10) / float64(days) / float64(pvbase)
	result, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", result), 64)
	if result < 0.5 {
		result = 0.5
	}
	return result, nil
}

// FindAllEditorFlowPaged 分页查询所有的编辑流量
// select e.username,e.realname,sum(f.pv_count) pv_count from baidu_url_editor e join baidu_url_flow f on e.page_id=f.page_id and time_span>='2020-04-05' and time_span<='2020-04-07
// ' group by e.username,e.realname order by pv_count desc  limit 30;
func (e *EditorTongji) FindAllEditorFlowPaged() error {
	euf, err := e.findAllEditorFlow()
	if err != nil {
		return err
	}
	e.Body.Data = append(e.Body.Data, euf)
	return nil
}

// FindAllEditorFlowPagedWithAvatar 查询编辑的流量并匹配微信头像
func (e *EditorTongji) FindAllEditorFlowPagedWithAvatar() error {
	euf, err := e.findAllEditorFlow()
	if err != nil {
		return err
	}
	if len(euf) == 0 {
		return nil
	}
	var usernames []string
	// log.Println("fff:", euf[0])
	for _, v := range euf {
		usernames = append(usernames, v.Realname)
	}
	wxs, err := FindWeixinStaffAvatarByUsername(usernames)
	if err != nil {
		return err
	}
	resluts := make(map[string]string)
	for _, val := range wxs {
		resluts[val.Username] = val.Avatar
	}
	for _, flow := range euf {
		flow.Avatar = resluts[flow.Realname]
		flow.Star, err = flow.GetInfluence(e.Body.StartDate, e.Body.EndDate)
		// log.Println("star:", flow.Star)
		if err != nil {
			return err
		}
	}
	e.Body.Data = append(e.Body.Data, euf)
	return nil
}
func (e *EditorTongji) findAllEditorFlow() ([]*EURLFlow, error) {
	var result []*EURLFlow
	if len(e.Body.StartDate) == 0 {
		e.Body.StartDate = util.FormatDate3(time.Now())
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = e.Body.StartDate
	}
	// if e.Body.MaxResults == 0 {
	// 	e.Body.MaxResults = 20
	// }

	var joins strings.Builder
	joins.WriteString("left join baidu_url_flow on baidu_url_flow.page_id=baidu_url_editor.page_id")
	if len(e.Body.UserName) > 0 {
		joins.WriteString(" and baidu_url_editor.username='" + e.Body.UserName + "'")
	}
	var total int
	// baidu_url_editor.username,baidu_url_editor.realname,sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time
	err := db.Table("baidu_url_editor").
		Select("baidu_url_editor.username,baidu_url_editor.realname,sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time").
		Joins(joins.String()).
		Group("baidu_url_editor.username,baidu_url_editor.realname").
		Where("baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", e.Body.StartDate, e.Body.EndDate).
		Count(&total).
		Order("pv_count desc,visitor_count desc").
		Offset(e.Body.StartIndex).
		Limit(e.Body.MaxResults).Find(&result).Error
	e.Body.Total = total
	return result, err
}

// FindEditorArticles 编辑的文章流量
func (e *EditorTongji) FindEditorArticles() error {
	if len(e.Body.UserName) == 0 {
		return errors.New("用户名[username:编辑的用户名,不是本名]不能为空")
	}
	now := time.Now()
	if len(e.Body.StartDate) == 0 {
		e.Body.StartDate = util.FormatDate3(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()))
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = util.FormatDate3(now)
	}
	if e.Body.MaxResults == 0 {
		e.Body.MaxResults = 20
	}
	var result []*EURLFlow
	var total int
	err := db.Table("baidu_url_editor").
		Select("baidu_url_flow.*,fz_manuscript.title,fz_manuscript.source").
		Joins("left join baidu_url_flow on baidu_url_flow.page_id=baidu_url_editor.page_id and baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", e.Body.StartDate, e.Body.EndDate).
		Joins("left join fz_manuscript on fz_manuscript.filename=baidu_url_flow.name").
		Where("baidu_url_editor.username=?", e.Body.UserName).
		Count(&total).
		Order("pv_count desc,visitor_count desc").
		Offset(e.Body.StartIndex).
		Limit(e.Body.MaxResults).Find(&result).Error
	e.Body.Data = append(e.Body.Data, result)
	e.Body.Total = total
	return err

}

// FindEditorDetails 查询编辑所有信息
func (e *EditorTongji) FindEditorDetails() error {
	if len(e.Body.UserName) == 0 {
		return errors.New("用户名[username:编辑的用户名,不是本名]不能为空")
	}
	now := time.Now()
	if len(e.Body.StartDate) == 0 {
		e.Body.StartDate = util.FormatDate3(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()))
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = util.FormatDate3(now)
	}
	// 查询编辑编辑过的文章及其流量列表
	// e.Body.Select = "sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time"
	err := e.FindAllEditorFlowPagedWithAvatar()
	if err != nil {
		return err
	}
	// 查询编辑的新老访客趋势
	err = e.FindEditorTrendVisitor()
	return err
}

// FindEditorTrendVisitor 新老用户访问趋势
func (e *EditorTongji) FindEditorTrendVisitor() error {
	if len(e.Body.UserName) == 0 {
		return errors.New("用户名[username:编辑的用户名,不是本名]不能为空")
	}
	now := time.Now()
	if len(e.Body.StartDate) == 0 {
		e.Body.StartDate = util.FormatDate3(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()))
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = util.FormatDate3(now)
	}
	vs := []string{"new", "old"}
	var result []*EURLFlow
	for _, visitor := range vs {
		// 分别查询新旧用户的流量
		flow := EURLFlow{}
		err := db.Table("baidu_url_flow").
			Select("visitor, sum(pv_count) as pv_count,sum(visitor_count) as visitor_count,sum(outward_count) as outward_count,sum(exit_count) as exit_count,round(sum(average_stay_time*pv_count)/sum(pv_count),0) as average_stay_time").
			Joins("join baidu_url_editor on username=? and baidu_url_flow.page_id=baidu_url_editor.page_id", e.Body.UserName).
			Where("baidu_url_flow.visitor=? and baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", visitor, e.Body.StartDate, e.Body.EndDate).
			Group("baidu_url_flow.visitor").Find(&flow).Error
		if err != nil {
			return err
		}
		if flow.PvCount != 0 {
			exitratio := fmt.Sprintf("%.2f", 100*float64(flow.ExitCount)/float64(flow.PvCount))
			flow.ExitRatio = exitratio
		}
		result = append(result, &flow)
		// log.Println("时长：", flow.AverageStayTime)
	}
	data := []interface{}{}
	fields := []string{"新旧访客比", "浏览量", "访客数", "跳出率", "平均访问时长"}
	// 新旧访客比
	newUvratio := float32(result[0].VisitorCount) / float32(result[0].VisitorCount+result[1].VisitorCount)
	oldUvratio := 1 - newUvratio
	data = append(data, []string{fmt.Sprintf("%.2f", 100*newUvratio), fmt.Sprintf("%.2f", 100*oldUvratio)})
	data = append(data, []int{result[0].PvCount, result[1].PvCount})
	data = append(data, []int{result[0].VisitorCount, result[1].VisitorCount})
	data = append(data, []string{result[0].ExitRatio, result[1].ExitRatio})
	data = append(data, []string{util.SecondsToTimesStr2(result[0].AverageStayTime), util.SecondsToTimesStr2(result[1].AverageStayTime)})
	e.Body.Data = append(e.Body.Data, []interface{}{data, fields})
	return nil
}

// FindFlowAndManuscriptNum 获取编辑的流量和稿件量
func (e *EditorTongji) FindFlowAndManuscriptNum() error {
	if len(e.Body.StartDate) == 0 {
		return errors.New("start_date不能为空")
	}
	if len(e.Body.EndDate) == 0 {
		e.Body.EndDate = e.Body.StartDate
	}

	var joins strings.Builder
	joins.WriteString("left join baidu_url_flow on baidu_url_flow.page_id=baidu_url_editor.page_id")
	if len(e.Body.UserName) > 0 {
		joins.WriteString(" and baidu_url_editor.username='" + e.Body.UserName + "'")
	}
	if len(e.Body.Metrics) == 0 {
		e.Body.Metrics = "baidu_url_editor.username,baidu_url_editor.realname,sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time"
	}
	// if len(e.Body.Metrics) == 0 {
	// 	e.Body.Metrics = "baidu_url_editor.username,baidu_url_editor.realname,count(fz_manuscript.editor),sum(baidu_url_flow.pv_count) as pv_count,sum(baidu_url_flow.visitor_count) as visitor_count,sum(baidu_url_flow.outward_count) as outward_count,sum(baidu_url_flow.exit_count) as exit_count,sum(baidu_url_flow.average_stay_time) as average_stay_time"
	// }

	// 查询流量
	var result []*EURLFlow
	var manuscript []*FzManuscriptCount
	var editors []*FzAdmin
	var err1, err2, err3 error
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		err1 = db.Table("baidu_url_editor").
			Select(e.Body.Metrics).
			Joins(joins.String()).
			// Joins("left join fz_manuscript on baidu_url_editor.username=fz_manuscript.editor and fz_manuscript.inserttime>=? and fz_manuscript.inserttime>=?", e.Body.StartDate, e.Body.EndDate).
			Group("baidu_url_editor.username,baidu_url_editor.realname").
			Where("baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", e.Body.StartDate, e.Body.EndDate).
			Order("pv_count desc,visitor_count desc").Find(&result).Error
		wg.Done()
	}()

	go func() {
		// 查询稿件数
		manuscript, err2 = CountEditorFzManuscriptFromLocal(map[string]interface{}{"start_date": e.Body.StartDate, "end_date": e.Body.EndDate})
		wg.Done()
	}()

	go func() {
		// 从远程查询当前有效的用户
		editors, err3 = FindEditorFromRemote()
		wg.Done()
	}()
	wg.Wait()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err1
	}
	if err3 != nil {
		return err1
	}
	editorMap := make(map[string]bool)
	for _, e := range editors {
		editorMap[e.Username+e.Realname] = true
	}
	// 删除编辑不存在的流量
	var indexs []int
	for i, flow := range result {
		if editorMap[flow.Username+flow.Realname] == false {
			indexs = append(indexs, i)
		}
	}
	size := len(indexs)
	for j := size - 1; j >= 0; j-- {
		result = append(result[:indexs[j]], result[indexs[j]+1:]...)
	}
	// 计算稿件量
	mmap := make(map[string]int)
	for _, m := range manuscript {
		mmap[m.Editor+m.Editorname] = m.Number
	}

	for _, flow := range result {
		flow.Star, _ = flow.GetInfluence(e.Body.StartDate, e.Body.EndDate)
		flow.Manuscript = mmap[flow.Username+flow.Realname]
	}
	e.Body.Data = append(e.Body.Data, result)
	return nil
}
