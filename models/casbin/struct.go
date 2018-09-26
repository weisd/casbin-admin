package casbin

// Permission Permission
type Permission struct {
	Name   string `json:"name" form:"name" query:"name"`
	Path   string `json:"path" form:"path" query:"path"`
	Method string `json:"method" form:"method" query:"method"`
	Origin string `json:"origin" form:"origin" query:"origin"`
}

// Parse Parse
func (p *Permission) Parse(fields []string) {
	if len(fields) < 4 {
		return
	}
	p.Name = fields[0]
	p.Path = fields[1]
	p.Method = fields[2]
	p.Origin = fields[3]
}

// RolePermission RolePermission
type RolePermission struct {
	Role       string `json:"role" form:"role" query:"role"`
	Permission string `json:"permission" form:"permission" query:"permission"`
}

// Parse ParseParse
func (p *RolePermission) Parse(fields []string) {
	if len(fields) < 2 {
		return
	}
	p.Role = fields[0]
	p.Permission = fields[1]
}

// UserRole UserRole
type UserRole struct {
	UID  string `json:"uid" form:"uid" query:"uid"`
	Role string `json:"role" form:"role" query:"role"`
}

// Request Request
type Request struct {
	UID    string `json:"uid" form:"uid" query:"uid"`
	Path   string `json:"path" form:"path" query:"path"`
	Method string `json:"method" form:"method" query:"method"`
	Origin string `json:"origin" form:"origin" query:"origin"`
}

// RolePermissionList RolePermissionList
type RolePermissionList struct {
	Role           string        `json:"role" form:"role" query:"role"`
	PermissionList []*Permission `json:"permission_list" form:"permission_list" query:"permission_list"`
}
