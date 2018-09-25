package admin

import "github.com/labstack/echo"
import "github.com/weisd/casbin-admin/router"

func init() {
	router.RegisterRouters("user", routers)
}

// routers routers
func routers(e *echo.Echo) {
	r := e.Group("user")

	r.GET("", Index)
	r.GET("/info", Info)

	r.POST("/login", Login)
	r.POST("/logout", Logout)
	r.POST("/add", Add)
	r.POST("/update", Update)
	r.POST("/update/pwd", UpdatePwd)

}
