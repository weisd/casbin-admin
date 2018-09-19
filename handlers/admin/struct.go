package admin

// AccountPasswd AccountPasswd
type AccountPasswd struct {
	Account string `json:"account" form:"account"`
	Passwd  string `json:"passwd" form:"passwd"`
}
