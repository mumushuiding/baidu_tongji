package model

import "github.com/mumushuiding/util"

// BaiduAPI 调用外部接口的对象
type BaiduAPI struct {
	Model
	URL       string `json:"url"`
	Operation string `json:"operation"` // POST GET
	Type      string `json:"type"`
	Params    Params `json:"params"`
	Par       string `json:"par" gorm:"size:1024"`
}

// Header 头
type Header struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	AccountType int    `json:"account_type"`
}

// Body body
type Body struct {
	SiteID     string `json:"site_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Metrics    string `json:"metrics"`
	Method     string `json:"method"`
	StartIndex int    `json:"start_index"`
	MaxResults int    `json:"max_results"`
}

// Params 参数
type Params struct {
	Header Header `json:"header"`
	Body   Body   `json:"body"`
}

// Params2Str 参数转字符串
func (a *BaiduAPI) Params2Str() string {
	str, _ := util.ToJSONStr(a.Params)
	return str
}

// GetParams 将Par字符串转换成Params对象
func (a *BaiduAPI) GetParams() {
	util.Str2Struct(a.Par, a.Params)
}

// ToString 转化成字符串
func (a *BaiduAPI) ToString() string {
	str, _ := util.ToJSONStr(a)
	return str
}

// FindAllAPI 使用map参数查询所有
func FindAllAPI(fields map[string]interface{}) ([]*BaiduAPI, error) {
	var result []*BaiduAPI
	return result, db.Where(fields).Find(&result).Error
}
