package admin

import "github.com/labstack/echo"
import "github.com/weisd/casbin-admin/router"

func init() {
	router.RegisterRouters("user", routers)
}

// routers routers
func routers(e *echo.Echo) {
	r := e.Group("user")
	r.POST("/login", Login)
	r.POST("/logout", Logout)
	r.POST("/add", Add)
	r.POST("/update/name", UpdateName)
	r.POST("/update/pwd", UpdatePwd)
}
