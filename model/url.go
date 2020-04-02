package model

// URL 保存不重复url地址
type URL struct {
	Model
	URL string `json:"url"`
}
