package jwt

// import (
// 	"fmt"
// 	"time"

// 	"github.com/jinzhu/gorm"
// 	_ "github.com/jinzhu/gorm/dialects/mysql"
// )

// // UserCAA UserCAA
// type UserCAA struct {
// 	ID        int64
// 	UID       int64
// 	Counter   int64
// 	Timeout   int64
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

// // MysqlStore MysqlStore
// type MysqlStore struct {
// 	DB *gorm.DB
// }

// // NewMysqlStore NewMysqlStore .eg root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
// func NewMysqlStore(dns string) (*MysqlStore, error) {
// 	db, err := gorm.Open("mysql", dns)
// 	if err != nil {
// 		return nil, fmt.Errorf("newGorm 连接数据库失败 %v", err)
// 	}

// 	db.SingularTable(true)

// 	if err := db.DB().Ping(); err != nil {
// 		return nil, err
// 	}

// 	if !db.HasTable(&UserCAA{}) {
// 		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARACTER SET=utf8").AutoMigrate(&UserCAA{})
// 	}

// 	s := &MysqlStore{}
// 	s.DB = db
// }

// // GetCounter GetCounter
// func (p *MysqlStore) GetCounter(uid int64) (int64, error) {
// 	info := &UserCAA{}
// 	g := p.DB.First(info, uid)
// 	if g.RecordNotFound() {
// 		return 0, nil
// 	}
// 	if g.Error != nil {
// 		return 0, g.Error
// 	}
// 	return info.Counter, nil
// }

// // IncrCounter IncrCounter
// func (p *MysqlStore) IncrCounter(uid int64, step int) error {

// 	return 0, nil
// }

// // GetTimeout GetTimeout
// func (p *MysqlStore) GetTimeout(uid int64) (int64, error) {
// 	return 0, nil
// }

// // SetTimeout SetTimeout
// func (p *MysqlStore) SetTimeout(uid int64, t int64) error {
// 	return 0, nil
// }
