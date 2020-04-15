package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
)

// 连接 129.0.99.111 的数据库
var dbNews *gorm.DB

// StartDBNews 启动数据库连接
func StartDBNews() {

	var err error
	dbNews, err = gorm.Open(conf.DbNewsType, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.DbNewsUser, conf.DbNewsPassword, conf.DbNewsHost, conf.DbNewsPort, conf.DbNewsName))
	if err != nil {
		log.Fatalf("数据库连接失败 err: %v", err)
	}
	log.Println("启动fznews_cms数据库连接！！")
	// 启用Logger，显示详细日志
	mode, _ := strconv.ParseBool(conf.DbLogMode)

	dbNews.LogMode(mode)
	dbNews.SingularTable(true) //全局设置表名不可以为复数形式
	idle, err := strconv.Atoi(conf.DbMaxIdleConns)
	dbNews.DB().SetMaxIdleConns(idle)
	open, err := strconv.Atoi(conf.DbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	dbNews.DB().SetMaxOpenConns(open)
	if err != nil {
		panic(err)
	}

}

// CloseDBNews 关闭数据库
func CloseDBNews() {
	defer dbNews.Close()
}
