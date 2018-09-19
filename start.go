package main

import (
	"fmt"

	"github.com/weisd/casbin-admin/models"

	gormadapter "github.com/casbin/gorm-adapter"
	_ "github.com/go-sql-driver/mysql" // need for gorm-adapter
	"github.com/weisd/casbin-admin/handlers"
	"github.com/weisd/casbin-admin/models/admin"
	"github.com/weisd/casbin-admin/models/casbin"
)

func start() {

	DB_USER := "root"
	DB_PWD := "root"
	DB_DATABASE := "casbin"
	DB_HOST := "localhost"

	models.Init(models.GormService{
		Host:   DB_HOST,
		DB:     DB_DATABASE,
		User:   DB_USER,
		Passwd: DB_PWD,
	})

	a := gormadapter.NewAdapter("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", DB_USER, DB_PWD, DB_HOST)) // Your driver and data source.
	casbin.Init(casbin.WithAdapter(a))

	handlers.Init()

	// init admin model
	admin.Init()
}
