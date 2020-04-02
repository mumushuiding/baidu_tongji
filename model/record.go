package model

// BaiduRecord 纪录
type BaiduRecord struct {
	Model
	Data string `json:"data" gorm:"size:1024"`
	Type string `json:"type"`
	Flag uint8  `json:"flag"` // 0失败，1成功
	Err  string `json:"err"`
}

// Save 直接保存
func (r *BaiduRecord) Save() error {
	return db.Create(r).Error
}

// FindLastRecord 查询第一条纪录
func FindLastRecord(fields map[string]interface{}) (*BaiduRecord, error) {
	var r BaiduRecord
	err := db.Where(fields).Last(&r).Error
	return &r, err
}
