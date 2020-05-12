package model

import (
	"github.com/mumushuiding/util"
)

// BaiduURLFlow url对应的流量
type BaiduURLFlow struct {
	Model
	TimeSpan string `json:"timeSpan,omitempty"`
	// 对就百度统计的url和id
	PageID          string `json:"pageId,omitempty"`
	Name            string `json:"name,omitempty"`
	PvCount         int    `json:"pv_count,omitempty"`
	VisitorCount    int    `json:"visitor_count,omitempty"`
	OutwardCount    int    `json:"outward_count,omitempty"`
	ExitCount       int    `json:"exit_count,omitempty"`
	AverageStayTime int    `json:"average_stay_time,omitempty"`
	Source          string `json:"source,omitempty"`
	Visitor         string `json:"visitor,omitempty"`
}

// FromMap 通过map赋值
func (uf *BaiduURLFlow) FromMap(fields map[string]interface{}) {
	// str, _ := util.ToJSONStr(fields)
	// log.Println("frommap:", str)
	uf.PvCount = fields["pv_count"].(int)
	uf.VisitorCount = fields["visitor_count"].(int)
	uf.OutwardCount = fields["outward_count"].(int)
	uf.ExitCount = fields["exit_count"].(int)
	if fields["average_stay_time"] != nil {
		uf.AverageStayTime = fields["average_stay_time"].(int)
	} else {
		uf.AverageStayTime = 0
	}

}

// ToString 转为字符串
func (uf *BaiduURLFlow) ToString() string {
	str, _ := util.ToJSONStr(uf)
	return str
}

// SaveOrUpdate 不存在就保存，存在就覆盖
func (uf *BaiduURLFlow) SaveOrUpdate() error {
	// return db.Create(uf).Error
	return db.Where(BaiduURLFlow{PageID: uf.PageID, TimeSpan: uf.TimeSpan}).Assign(uf).FirstOrCreate(uf).Error
}

// Save 不存在就保存，存在就覆盖
func (uf *BaiduURLFlow) Save() error {
	return db.Create(uf).Error
}

// FindLast 最新纪录
func (uf *BaiduURLFlow) FindLast() error {
	return db.Table("baidu_url_flow").Last(uf).Error
}
