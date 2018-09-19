package admin

import "github.com/weisd/casbin-admin/models"

// Init Init model 初始化后执行
func Init() {
	if !models.DB.HasTable(&CasbinAdmin{}) {
		models.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARACTER SET=utf8").AutoMigrate(&CasbinAdmin{})
	}
}
