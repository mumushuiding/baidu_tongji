package model

import (
	"errors"
	"strings"
)

// WeixinStaff 编辑对应的微账号信息
type WeixinStaff struct {
	ID       int    `json:"id"`
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
}

// FindWeixinStaffAvatarByUsername 根据用户名返回用户的头像
// 用户名重复的返回其中id小的
// select username,avatar from weixin_staff a where not exists(select 1 from weixin_staff b where a.username=b.username and b.id<a.id) and username in('admin1','王月玲');
func FindWeixinStaffAvatarByUsername(usernames []string) ([]*WeixinStaff, error) {
	var result []*WeixinStaff
	var build strings.Builder
	if len(usernames) == 0 {
		return nil, errors.New("用户名数组不能为空")
	}
	for _, v := range usernames {
		build.WriteString("'" + v + "',")
	}
	u := build.String()
	where := "select username,avatar from weixin_staff where username in(" + u[0:len(build.String())-1] + ")"
	err := db.Raw(where).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
