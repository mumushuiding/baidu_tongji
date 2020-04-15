package controller

import (
	"net/http"

	"github.com/mumushuiding/baidu_tongji/conmgr"
	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/util"
)

// GetBaiduDataByTimeSpan 根据时间跨度抓取百度统计数据
// 参数: dates []string
func GetBaiduDataByTimeSpan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法,参数开始日期startDate: 20200102,结束日期endDate: 20200103")
		return
	}
	// 如果body参数有URL编码的话会报错
	body, err := util.Body2Map(r)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	if body["startDate"] == nil {
		util.ResponseErr(w, "startDate不能为空，参数如: startDate:20200103")
		return
	}
	if body["endDate"] == nil {
		util.ResponseErr(w, "startDate不能为空，参数如: startDate:20200103")
		return
	}
	startDate, ok := body["startDate"].(string)
	if !ok {
		util.ResponseErr(w, "startDate 类型必须是字符串")
		return
	}
	endDate, ok := body["endDate"].(string)
	if !ok {
		util.ResponseErr(w, "endDate 类型必须是字符串")
		return
	}
	// dates, ok := body["dates"].([]interface{})
	// if !ok {
	// 	util.ResponseErr(w, "dates类型必须是数组")
	// 	return
	// }
	// var res []string
	// for _, d := range dates {
	// 	res = append(res, d.(string))
	// }
	// if len(res) == 0 {
	// 	return
	// }
	err = conmgr.GetBaiduDataByTimeSpan(startDate, endDate)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseData(w, "成功发起查询，查询中。。")

}

// GetTongjiData 获取统计数据
func GetTongjiData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法")
		return
	}
	var body model.EditorTongji
	err := util.Body2Struct(r, &body)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	// log.Println("body:", body.ToString())
	if len(body.Body.StartDate) > 0 {
		_, err := util.ParseDate3(body.Body.StartDate)
		if err != nil {
			util.ResponseErr(w, err)
			return
		}
	}
	if len(body.Body.EndDate) > 0 {
		_, err := util.ParseDate3(body.Body.EndDate)
		if err != nil {
			util.ResponseErr(w, err)
			return
		}
	}
	f, err := GetRoute(body.Body.Method)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	result, err := f(&body)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseData(w, result)

}

// GetRemoteFzManuscriptNotHave 远程拉取稿件数据
func GetRemoteFzManuscriptNotHave(w http.ResponseWriter, r *http.Request) {
	conmgr.Conmgr.GetRemoteFzManuscriptNotHave()
	util.ResponseData(w, "已经开始从远程拉取数据，失败的纪录会保存在,typename为【recordgetRemoteFzManuscriptNotHave】")
}
