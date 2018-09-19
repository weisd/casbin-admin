package router

import "github.com/labstack/echo"

// RouterFunc RouterFunc
type RouterFunc func(e *echo.Echo)

// Routers router store
var Routers = map[string]RouterFunc{}

// Init Init
func Init(e *echo.Echo) {
	for i := range Routers {
		Routers[i](e)
	}
}

// RegisterRouters RegisterRouters
func RegisterRouters(name string, r RouterFunc) {
	Routers[name] = r
}
