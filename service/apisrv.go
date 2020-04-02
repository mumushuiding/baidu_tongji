package service

import (
	"github.com/mumushuiding/baidu_tongji/model"
)

// GetTaskFromAPI 从api表获取需要执行的查询任务
func GetTaskFromAPI() ([]*model.BaiduAPI, error) {
	return model.FindAllAPI(make(map[string]interface{}))
}
