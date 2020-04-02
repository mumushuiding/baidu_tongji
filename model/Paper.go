package model

// FzManuscript 稿件
type FzManuscript struct {
	Filename   string `json:"filename"`
	Editor     string `json:"editor"`
	Editorname string `json:"editorname"`
	// 对应百度统计的PageId
	PageID string `json:"pageId"`
}

// Transform2URLEditor 转成BaiduURLEditor对象
func (p *FzManuscript) Transform2URLEditor() BaiduURLEditor {
	var r BaiduURLEditor
	r.PageID = p.PageID
	r.Realname = p.Editorname
	r.Username = p.Editor
	return r
}

// FindPaperFirst 查询第一条稿件
func FindPaperFirst(fields map[string]interface{}) (*FzManuscript, error) {
	var p FzManuscript
	err := db.Where(fields).First(&p).Error
	return &p, err
}
