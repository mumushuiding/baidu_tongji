package model

import (
	"errors"
	"strings"
	"time"

	"github.com/mumushuiding/util"
)

// FzManuscript 稿件
type FzManuscript struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Suptitle   string    `json:"suptitle"`
	Subtitle   string    `json:"subtitle"`
	Shorttitle string    `json:"shorttitle"`
	Publictime time.Time `json:"publictime"`
	Writer     string    `json:"writer"`
	Source     string    `json:"source"`
	Type       int       `json:"type"`
	State      int       `json:"state"`
	Tags       string    `json:"tags"`
	Channelid  string    `json:"channelid"`
	Click      int       `json:"click"`
	Filename   string    `json:"filename"`
	Editor     string    `json:"editor"`
	Editorname string    `json:"editorname"`

	// 对应百度统计的PageId,一个filename为有多个pageID，因为来源不同，所以这个字段，没有统计作用
	PageID string `json:"pageId"`
}

// FzManuscriptCount 稿件数
type FzManuscriptCount struct {
	Editor     string `json:"editor"`
	Editorname string `json:"editorname"`
	Number     int    `json:"number"`
}

// FBody FBody
// type FBody struct {
// 	Data []interface{} `json:data,omitempty`
// }

// FzManuscriptNotFound 本地数据库不存在的稿件
type FzManuscriptNotFound struct {
	Filename string
	PageID   string
}

// SaveOrUpdate 存在就覆盖
func (p *FzManuscript) SaveOrUpdate() error {

	return db.Where(FzManuscript{Filename: p.Filename}).Assign(*p).Omit("page_id").FirstOrCreate(p).Error
}

// CountEditorFzManuscriptFromLocal 从本地查询编辑对应的稿件数
func CountEditorFzManuscriptFromLocal(fields map[string]interface{}) ([]*FzManuscriptCount, error) {
	var where strings.Builder
	if fields["start_date"] != nil {
		where.WriteString(" and fz_manuscript.inserttime>='" + fields["start_date"].(string) + "'")
	}
	if fields["end_date"] != nil {
		where.WriteString(" and fz_manuscript.inserttime<='" + fields["end_date"].(string) + "'")
	}
	if fields["editor"] != nil {
		where.WriteString(" and fz_manuscript.editor='" + fields["editor"].(string) + "'")
	}
	where.WriteString(" and state!=-1")
	var w string
	if len(where.String()) > 0 {
		w = where.String()[4:]
	}
	var result []*FzManuscriptCount
	err := db.Select("editor,editorname,count(id) as number").Table("fz_manuscript").
		Where(w).
		Group("editor,editorname").
		Order("number desc").
		Find(&result).
		Error
	return result, err
}

// CountFzManuscriptFromLocal 条数
func CountFzManuscriptFromLocal(fields map[string]interface{}) (int, error) {
	var count int
	var where strings.Builder
	if fields["start_date"] != nil {
		where.WriteString("and inserttime>='" + fields["start_date"].(string) + "'")
	}
	if fields["end_date"] != nil {
		where.WriteString("and inserttime<='" + fields["end_date"].(string) + "'")
	}
	if fields["editor"] != nil {
		where.WriteString("and editor='" + fields["editor"].(string) + "'")
	}
	var w string
	if len(where.String()) > 0 {
		w = where.String()[4:]
	}
	err := db.Table("fz_manuscript").Where(w).Count(&count).Error
	return count, err
}

// FindFzManuscriptFromLocal 本地查询
func FindFzManuscriptFromLocal(fields map[string]interface{}) ([]*FzManuscript, int, error) {
	if len(fields) == 0 {
		return nil, 0, errors.New("查询参数不能全为空")
	}
	if fields["max_results"] == nil {
		fields["max_results"] = 100
	}
	if fields["order"] == nil {
		fields["order"] = "inserttime desc"
	}
	if fields["start_index"] == nil {
		fields["start_index"] = 0
	}
	if fields["where"] == nil {
		fields["where"] = ""
	}
	var r []*FzManuscript
	var count int
	err := db.
		Table("fz_manuscript").
		Where(fields["where"]).
		Count(&count).
		Order(fields["order"]).
		Offset(fields["start_index"]).
		Limit(fields["max_results"]).
		Find(&r).Error
	return r, count, err
}

// select title,editor,name,sum(pv_count) as pv_count from baidu_url_flow b join fz_manuscript f on (f.filename=b.name and time_span>='2020-04-06' and time_span<='2020-05-01')  group by title,editor,name order by pv_count desc limit 10;

// FindFzManuscripFlow 多表查询文章的流量和编辑的头像
func FindFzManuscripFlow(fields map[string]interface{}) ([]*EURLFlow, error) {
	if len(fields) == 0 {
		return nil, errors.New("查询参数不能全为空")
	}
	if fields["max_results"] == nil {
		fields["max_results"] = 50
	}
	if fields["order"] == nil {
		fields["order"] = "pv_count desc,visitor_count desc"
	}
	if fields["start_index"] == nil {
		fields["start_index"] = 0
	}
	if fields["start_date"] == nil {
		fields["start_date"] = util.FormatDate3(time.Now())
	}
	if fields["end_date"] == nil {
		fields["end_date"] = fields["start_date"]
	}
	var r []*EURLFlow
	err := db.Table("baidu_url_flow").
		Select("baidu_url_flow.*,fz_manuscript.source,fz_manuscript.title,fz_manuscript.editor as username,fz_manuscript.editorname as realname,fz_manuscript.filename").
		Joins("join fz_manuscript on fz_manuscript.filename=baidu_url_flow.name").
		Where("baidu_url_flow.time_span>=? and baidu_url_flow.time_span<=?", fields["start_date"], fields["end_date"]).
		Order(fields["order"]).
		Offset(fields["start_index"]).
		Limit(fields["max_results"]).Find(&r).Error
	return r, err

}

// FindFzManuscriptByURLFromLocalDB 本地查询
func FindFzManuscriptByURLFromLocalDB(fields map[string]interface{}) (*FzManuscript, error) {
	var f FzManuscript
	err := db.Where(fields).First(&f).Error
	return &f, err
}

// FindFirstFzManuscriptFromDBNews 从DBNews数据库查询稿件，返回第一条
func FindFirstFzManuscriptFromDBNews(fields map[string]interface{}) (*FzManuscript, error) {
	var f FzManuscript
	err := dbNews.Where(fields).First(&f).Error
	return &f, err
}

// FindFzManuscriptFromDBNewsByFilenames 根据filename批量查询稿件
func FindFzManuscriptFromDBNewsByFilenames(filenames []string) ([]*FzManuscript, error) {
	var r []*FzManuscript
	err := dbNews.Where("filename in (?)", filenames).Find(&r).Error
	return r, err
}

// ToString ToString
func (p *FzManuscript) ToString() string {
	str, _ := util.ToJSONStr(p)
	return str
}

// Transform2URLEditor 转成BaiduURLEditor对象
func (p *FzManuscript) Transform2URLEditor() BaiduURLEditor {
	var r BaiduURLEditor
	r.PageID = p.PageID
	r.Realname = p.Editorname
	r.Username = p.Editor
	return r
}

// FindFzManuscriptFirstFromDBNews 远程查询第一条稿件
func FindFzManuscriptFirstFromDBNews(fields map[string]interface{}) (*FzManuscript, error) {
	var p FzManuscript
	err := dbNews.Where(fields).First(&p).Error
	return &p, err
}

// FindFzManuscriptFromDBNews 远程查询
func FindFzManuscriptFromDBNews(fields map[string]interface{}) ([]*FzManuscript, int, error) {
	if len(fields) == 0 {
		return nil, 0, errors.New("查询参数不能全为空")
	}
	if fields["max_results"] == nil {
		fields["max_results"] = 100
	}
	if fields["order"] == nil {
		fields["order"] = "inserttime desc"
	}
	if fields["start_index"] == nil {
		fields["start_index"] = 0
	}
	if fields["where"] == nil {
		fields["where"] = ""
	}
	var r []*FzManuscript
	var count int
	err := dbNews.
		Table("fz_manuscript").
		Where(fields["where"]).
		Count(&count).
		Order(fields["order"]).
		Offset(fields["start_index"]).
		Limit(fields["max_results"]).
		Find(&r).Error
	return r, count, err
}
