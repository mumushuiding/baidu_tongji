package model

// FzAdmin FzAdmin
type FzAdmin struct {
	Username string `json:"username"`
	Realname string `json:"realname"`
	// 1为编辑
	Department int `json:"department"`
}

// FindEditorFromRemote 从远程查询编辑
func FindEditorFromRemote() ([]*FzAdmin, error) {
	var result []*FzAdmin
	err := dbNews.Table("fz_admin").Where("department=1").Find(&result).Error
	return result, err
}
