package admin

// AccountPasswd AccountPasswd
type AccountPasswd struct {
	Account string `json:"account" form:"account"`
	Passwd  string `json:"passwd" form:"passwd"`
}

// LoginResp LoginResp
type LoginResp struct {
	Token string `json:"token"`
}

// ID ID
type ID struct {
	ID int64 `json:"id" form:"id" query:"id" validate:"gt=0"`
}

// CasbinAdminCreate CasbinAdminCreate
type CasbinAdminCreate struct {
	Name   string `json:"name" form:"name" query:"name" validate:"required,min=3,max=10"`
	Email  string `json:"email" form:"email" query:"email"  validate:"omitempty,email"`
	Phone  string `json:"phone" form:"phone" query:"phone"`
	Passwd string `json:"passwd" form:"passwd" query:"passwd" validate:"required,min=6"`
	Status int32  `json:"status" form:"status" query:"status" `
}

// CasbinAdminUpdate CasbinAdminUpdate
type CasbinAdminUpdate struct {
	ID     int64  `json:"id" form:"id" query:"id" validate:"gt=0"`
	Name   string `json:"name" form:"name" query:"name" validate:"required,min=3,max=10"`
	Email  string `json:"email" form:"email" query:"email"  validate:"email"`
	Phone  string `json:"phone" form:"phone" query:"phone"`
	Status int32  `json:"status" form:"status"  query:"status" `
}

// CasbinAdminSearchArgs CasbinAdminSearchArgs
type CasbinAdminSearchArgs struct {
	Querys []string `json:"querys" form:"querys[]" query:"querys[]"`
	Limit  int      `json:"limit" form:"limit" query:"limit"`
	Order  string   `json:"order" form:"order" query:"order"`
}
