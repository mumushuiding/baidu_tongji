package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/baidu_tongji/config"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

// Model 其它数据结构的公共部分
type Model struct {
	ID         int       `gorm:"primary_key" json:"id,omitempty"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime,omitempty"`
}

// 配置
var conf = *config.Config

// SetupDB 初始化一个db连接
func SetupDB() {
	var err error
	db, err = gorm.Open(conf.DbType, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName))
	if err != nil {
		log.Fatalf("数据库连接失败 err: %v", err)
	}
	log.Println("启动数据库连接！！")
	// 启用Logger，显示详细日志
	mode, _ := strconv.ParseBool(conf.DbLogMode)

	db.LogMode(mode)

	idle, err := strconv.Atoi(conf.DbMaxIdleConns)
	db.DB().SetMaxIdleConns(idle)
	// 防止连接时间超过数据库上限，导致连接失效
	db.DB().SetConnMaxLifetime(time.Minute * 10)
	open, err := strconv.Atoi(conf.DbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(open)
	if err != nil {
		panic(err)
	}

	db.SingularTable(true) //全局设置表名不可以为复数形式
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(&BaiduAPI{}).AutoMigrate(&BaiduRecord{}).AutoMigrate(&BaiduURLFlow{}).AutoMigrate(&BaiduURLEditor{}).AutoMigrate(&FzManuscript{})
	db.Model(&BaiduURLEditor{}).AddUniqueIndex("pageid", "page_id")
	db.Model(&BaiduURLEditor{}).AddUniqueIndex("pageid_username", "page_id", "username")
	db.Model(&BaiduURLEditor{}).AddIndex("username", "username")
	db.Model(&BaiduURLFlow{}).AddUniqueIndex("pageid_timespan_source_visitor", "page_id", "time_span", "source", "visitor")
	db.Model(&BaiduURLFlow{}).AddIndex("name", "name")
	db.Model(&BaiduURLFlow{}).AddIndex("timespan", "time_span")
	db.Model(&FzManuscript{}).AddIndex("filename_index", "filename")
	db.Model(&FzManuscript{}).AddIndex("editor_index", "editor")
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

// GetDB getdb
func GetDB() *gorm.DB {
	return db
}

// GetTx GetTx
func GetTx() *gorm.DB {
	return db.Begin()
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		// if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
		// 	if modifyTimeField.IsBlank {
		// 		modifyTimeField.Set(nowTime)
		// 	}
		// }
	}
}
