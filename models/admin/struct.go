package admin

import "time"

// CasbinAdmin CasbinAdmin
type CasbinAdmin struct {
	ID        int64     `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	Email     string    `json:"email" form:"email"`
	Phone     string    `json:"phone" form:"phone"`
	Passwd    string    `json:"passwd" form:"passwd"`
	Salt      string    `json:"salt" form:"salt"`
	Status    int32     `json:"status" form:"status"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

// SafeInfo SafeInfo
func (p *CasbinAdmin) SafeInfo() *CasbinAdmin {
	p.Passwd = ""
	p.Salt = ""
	return p
}
