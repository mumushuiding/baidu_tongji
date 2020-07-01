package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/util"
)

// GetDatasOfTargetAPI 从指定接口获取数据
func GetDatasOfTargetAPI(url string, params string) (*model.BaiduData, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(params))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var baidudata model.BaiduData
	err = json.NewDecoder(resp.Body).Decode(&baidudata)
	if err != nil {
		return nil, err
	}
	return &baidudata, nil
}

// GetTotalNumOfTargetAPI 获取指定查询API总条数
func GetTotalNumOfTargetAPI(url string, params string) (int, error) {
	baidudata, err := GetDatasOfTargetAPI(url, params)
	if err != nil {
		s, _ := util.ToJSONStr(baidudata)
		return 0, fmt.Errorf(`{"data":%s,"err":%s}`, s, err.Error())
	}
	return baidudata.GetTotalNums(), nil
}

// GetDatesToFind 获取本次需要查询的日期数组
func GetDatesToFind(typename string) ([]string, error) {
	// 查询上次查询的日期 yyyymmdd,查不到就返回今天
	// 如果小于今天就添加上今天
	var result model.BaiduURLFlow
	err := result.FindLast()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	endDate := today.Add(time.Hour * (-24))
	if err == gorm.ErrRecordNotFound {
		return []string{util.FormatDate1(endDate)}, nil
	}
	var dates []string
	log.Println("最后更新时间:", result.TimeSpan)
	timespan, _ := util.ParseDate3(result.TimeSpan)
	timespan = time.Date(timespan.Year(), timespan.Month(), timespan.Day(), 0, 0, 0, 0, timespan.Location())

	// sub := starDate.Day() - timespan.Day()
	sub, _ := util.DateSubReturnDays(timespan, endDate)
	if sub > 0 {
		for i := 1; i <= sub; i++ {
			dates = append(dates, util.FormatDate1(time.Date(timespan.Year(), timespan.Month(), timespan.Day()+i, 0, 0, 0, 0, timespan.Location())))
		}
	}
	return dates, nil
}

// SaveURLWebFlow 保存网址流量
func SaveURLWebFlow(urlflow *model.BaiduURLFlow) error {
	// 存在就覆盖
	return urlflow.SaveOrUpdate()
}
