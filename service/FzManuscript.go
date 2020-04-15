package service

import "github.com/mumushuiding/baidu_tongji/model"

// FindFzManuscriptFromLocal 本地查询
func FindFzManuscriptFromLocal(fields map[string]interface{}) ([]*model.FzManuscript, int, error) {
	return model.FindFzManuscriptFromLocal(fields)
}

// FindFzManuscriptByURLFromLocalDB 本地查询
func FindFzManuscriptByURLFromLocalDB(url string) (*model.FzManuscript, error) {
	return model.FindFzManuscriptByURLFromLocalDB(map[string]interface{}{"filename": url})
}

// FindFzManuscriptFromDBNews 远程查询
func FindFzManuscriptFromDBNews(fields map[string]interface{}) ([]*model.FzManuscript, int, error) {
	return model.FindFzManuscriptFromDBNews(fields)
}

// FindFzManuscriptByURLFromDBNews 远程查询
func FindFzManuscriptByURLFromDBNews(url string) (*model.FzManuscript, error) {
	return model.FindFzManuscriptFirstFromDBNews(map[string]interface{}{"filename": url})
}

// FindFzManuscriptFromDBNewsByFilenames 根据filenames批量查询
func FindFzManuscriptFromDBNewsByFilenames(filenames []string) ([]*model.FzManuscript, error) {
	return model.FindFzManuscriptFromDBNewsByFilenames(filenames)
}
