package model

import (
	"errors"
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

// FBody FBody
// type FBody struct {
// 	Data []interface{} `json:data,omitempty`
// }

// FzManuscriptNotFound 本地数据库不存在的稿件
type FzManuscriptNotFound struct {
	Filename string
	PageID   string
}

// FindFzManuscriptFromLocal 远程查询
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

// SaveOrUpdate 存在就覆盖
func (p *FzManuscript) SaveOrUpdate() error {

	return db.Where(FzManuscript{Filename: p.Filename}).Assign(&p).FirstOrCreate(p).Error
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

// CountFzManuscriptFromDBNews 条数
func CountFzManuscriptFromDBNews(fields map[string]interface{}) (int, error) {
	var count int
	err := dbNews.Where(fields).Count(&count).Error
	return count, err
}