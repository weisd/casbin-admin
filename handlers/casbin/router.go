package casbin

import "github.com/labstack/echo"
import "github.com/weisd/casbin-admin/router"

func init() {
	router.RegisterRouters("casbin", routers)
}

// routers routers
func routers(e *echo.Echo) {
	c := e.Group("casbin")

	c.GET("", Index)

	c.POST("/model", GetModel)
	c.POST("/enforce", Enforce)

	p := c.Group("/permission")
	p.POST("/list", PermissionList)
	p.POST("/add", PermissionAdd)
	p.POST("/del", PermissionDel)

	r := c.Group("/role")
	r.POST("/list", RoleList)
	r.POST("/permission/list", RolePermissions)
	r.POST("/permission/add", RolePermissionAdd)
	r.POST("/permission/del", RolePermissionDel)

	u := c.Group("/user")
	u.POST("/role/list", UserRoles)
	u.POST("/role/add", UserRoleAdd)
	u.POST("/role/del", UserRoleDel)
	u.POST("/permission/list", UserPermissions)
}
