package model

import "github.com/mumushuiding/util"

// BaiduURLEditor 网址对应编辑
type BaiduURLEditor struct {
	Model
	PageID   string `json:"pageId"`
	Username string `json:"username,omitempty"`
	Realname string `json:"realname,omitempty"`
}

// SaveOrUpdate 存在就覆盖
func (e *BaiduURLEditor) SaveOrUpdate() error {
	return db.Where(BaiduURLEditor{PageID: e.PageID}).Assign(e).FirstOrCreate(e).Error
}

// ToString ToString
func (e *BaiduURLEditor) ToString() string {
	str, _ := util.ToJSONStr(e)
	return str
}
