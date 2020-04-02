package service

import "github.com/mumushuiding/baidu_tongji/model"

// FindFzManuscriptByURL 根据网址查询稿件
func FindFzManuscriptByURL(url string) (*model.FzManuscript, error) {
	return model.FindPaperFirst(map[string]interface{}{"filename": url})
}
