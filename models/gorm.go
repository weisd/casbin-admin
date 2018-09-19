package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// GormService GormService
type GormService struct {
	Host   string
	DB     string
	User   string
	Passwd string
}

// DB DBDB
var DB *gorm.DB

// Init Init
func Init(conf GormService) {
	var err error
	DB, err = newGorm(conf)
	if err != nil {
		panic(err)
	}
}

func newGorm(conf GormService) (*gorm.DB, error) {

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Passwd, conf.Host, conf.DB))
	if err != nil {
		return nil, fmt.Errorf("newGorm 连接数据库失败 %v", err)
	}

	db.SingularTable(true)

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
