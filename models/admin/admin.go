package admin

import (
	"github.com/weisd/casbin-admin/models"
)

// CheckPasswd CheckPasswd
func CheckPasswd(info *CasbinAdmin, passwd string) bool {

	if len(info.Passwd) > 0 && HashPassword(passwd, info.Salt) == info.Passwd {
		return true
	}

	return false
}

// InfoByAccount InfoByAccountInfoByAccount
func InfoByAccount(account string) (*CasbinAdmin, error) {
	switch {
	case IsPhone(account):
		return InfoByPhone(account)
	case IsEmail(account):
		return InfoByEmail(account)
	default:
		return InfoByName(account)
	}
}

// InfoByID InfoByID
func InfoByID(id int64) (*CasbinAdmin, error) {
	info := &CasbinAdmin{}
	g := models.DB.First(info, id)
	if g.RecordNotFound() {
		return nil, nil
	}
	if g.Error != nil {
		return nil, g.Error
	}

	return info, nil
}

// InfoByName InfoByName
func InfoByName(name string) (*CasbinAdmin, error) {
	info := &CasbinAdmin{}
	g := models.DB.Where("name = ?", name).First(info)
	if g.RecordNotFound() {
		return nil, nil
	}
	if g.Error != nil {
		return nil, g.Error
	} else if g.RecordNotFound() {
		return nil, nil
	}

	return info, nil
}

// InfoByPhone InfoByPhone
func InfoByPhone(phone string) (*CasbinAdmin, error) {
	info := &CasbinAdmin{}
	g := models.DB.Where("phone = ?", phone).First(info)
	if g.RecordNotFound() {
		return nil, nil
	}
	if g.Error != nil {
		return nil, g.Error
	} else if g.RecordNotFound() {
		return nil, nil
	}

	return info, nil
}

// InfoByEmail InfoByEmail
func InfoByEmail(email string) (*CasbinAdmin, error) {
	info := &CasbinAdmin{}
	g := models.DB.Where("email = ?", email).First(info)
	if g.RecordNotFound() {
		return nil, nil
	}
	if g.Error != nil {
		return nil, g.Error
	} else if g.RecordNotFound() {
		return nil, nil
	}

	return info, nil
}

// Create Create
func Create(info *CasbinAdmin) error {

	info.Salt = RandString(4)
	info.Passwd = HashPassword(info.Passwd, info.Salt)
	info.Status = 1

	g := models.DB.Create(info)
	if g.Error != nil {
		return g.Error
	}

	return nil
}

// UpdateName UpdateName
func UpdateName(id int64, name string) error {
	g := models.DB.Model(&CasbinAdmin{}).Where("id = ?", id).Update("name", name)
	if g.Error != nil {
		return g.Error
	}

	return nil
}

// UpdatePwd UpdatePwd
func UpdatePwd(id int64, pwd string) error {

	salt := RandString(4)
	pwd = HashPassword(pwd, salt)

	g := models.DB.Model(&CasbinAdmin{}).Where("id = ?", id).Update(CasbinAdmin{Passwd: pwd, Salt: salt})
	if g.Error != nil {
		return g.Error
	}

	return nil
}

// UpdateStatus UpdateStatus
func UpdateStatus(id int64, status int32) error {
	g := models.DB.Model(&CasbinAdmin{}).Where("id = ?", id).Update("status", status)
	if g.Error != nil {
		return g.Error
	}

	return nil
}
