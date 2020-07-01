package model

import (
	"time"

	"github.com/mumushuiding/util"
)

// BaiduData 百度统计返回数据
type BaiduData struct {
	Header BHeader `json:"header"`
	Body   BBody   `json:"body"`
}

// BHeader 头
type BHeader struct {
	Desc     string        `json:"desc,omitempty"`
	Failures []interface{} `json:"failures,omitempty"`
	Oprs     int           `json:"oprs,omitempty"`
	Succ     int           `json:"succ,omitempty"`
	Oprtime  int           `json:"oprtime,omitempty"`
	Quota    int           `json:"quota,omitempty"`
	Rquota   int           `json:"rquota,omitempty"`
	Status   int           `json:"status,omitempty"`
}

// BItems 页面连接
type BItems struct {
	PageID string `json:"pageId"`
	Name   string `json:"name"`
}

// BResult 结果
type BResult struct {
	Total    int               `json:"total"`
	Items    [][][]interface{} `json:"items"`
	TimeSpan []string          `json:"timeSpan"`
	Sum      [][]int           `json:"sum"`
	Offset   int               `json:"offset"`
	PageSum  [][]int           `json:"pageSum"`
	Fields   []string          `json:"fields"`
}

// BData BData
type BData struct {
	Result BResult `json:"result"`
}

// BBody BBody
type BBody struct {
	Data    []BData `json:"data"`
	Source  string  `json:"source"`
	Visitor string  `json:"visitor"`
}

// IfVisitSucess 判断是否访问成功
// 已知 Status 为 0 时成功
// status 为 2 时没有权限
func (b *BaiduData) IfVisitSucess() bool {
	if b.Header.Status == 0 {
		return true
	}
	return false
}

// GetTotalNums 获取总条数
func (b *BaiduData) GetTotalNums() int {
	if !b.IfVisitSucess() {
		return 0
	}

	return b.Body.Data[0].Result.Total
}

// ToString 转换成字符串
func (b *BaiduData) ToString() string {
	str, _ := util.ToJSONStr(b)
	return str
}

// GetBItems 获取百度统计中URL链接名对象,
// 如：[{"pageId": "15224687311766939350","name": "http://news.fznews.com.cn/dsxw/20200228/5e58c26f58664.shtml?from=m"}]
func (b *BaiduData) GetBItems() []BItems {
	var result []BItems
	urls := b.Body.Data[0].Result.Items[0]
	for _, item := range urls {
		var b BItems
		url := item[0].(map[string]interface{})
		b.PageID = url["pageId"].(string)
		b.Name = url["name"].(string)
		result = append(result, b)
	}
	return result
}

// Transform2URLFlow 提取受访页面流量
func (b *BaiduData) Transform2URLFlow() []*BaiduURLFlow {
	items := b.Body.Data[0].Result.Items
	urls := items[0]
	flows := items[1]
	fields := b.Body.Data[0].Result.Fields
	var result []*BaiduURLFlow
	for i, item := range urls {
		// 流量
		if len(flows[i]) == 0 {
			continue
		}
		var flow = make(map[string]interface{})
		for j, field := range fields {
			// 0 为 visit_page_title 不需要
			if j == 0 {
				continue
			}
			switch flows[i][j-1].(type) {
			case float64:
				flow[field] = int(flows[i][j-1].(float64))
				break
			}

		}
		var uf = &BaiduURLFlow{}
		uf.Source = b.Body.Source
		uf.Visitor = b.Body.Visitor
		uf.TimeSpan = b.Body.Data[0].Result.TimeSpan[0]
		uf.CreateTime = time.Now()
		// if len(flow) == 0 {
		// 	continue
		// }
		uf.FromMap(flow)
		url := item[0].(map[string]interface{})
		uf.PageID = url["pageId"].(string)
		uf.Name = url["name"].(string)

		// str, _ := util.ToJSONStr(uf)
		// log.Println("uf:", str)
		result = append(result, uf)

	}
	return result
}
