package service

import (
	"github.com/mumushuiding/baidu_tongji/model"
	"github.com/mumushuiding/util"
)

// HandleErrRecord 处理错误
func HandleErrRecord(e *model.EditorTongji) error {
	result, err := model.FindRecordByTypeAndErr("存储FzManuscript失败", "Error 1054: Unknown column 'page_id' in 'field list'")
	if err != nil {
		return err
	}
	for _, r := range result {
		var fm model.FzManuscript
		err := util.Str2Struct(r.Data, &fm)
		if err != nil {
			return err
		}
		err = fm.SaveOrUpdate()
		if err != nil {
			return err
		}
		err = r.Del()
		if err != nil {
			return err
		}
	}
	return nil
}
