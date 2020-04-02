package model

// EditorTongji 用户统计数据
type EditorTongji struct {
	Header BHeader `json:"header"`
	Body   EBody   `json:"body"`
}

// EData EData
type EData struct {
	Username string       `json:"username"`
	Realname string       `json:"realname"`
	TimeSpan []string     `json:"timeSpan,omitempty"`
	URLFlow  BaiduURLFlow `json:"urlFlow,omitempty"`
}

// EBody EBody
type EBody struct {
	Data       []EData `json:"data"`
	Total      int     `json:"total"`
	StartIndex int     `json:"start_index"`
}
