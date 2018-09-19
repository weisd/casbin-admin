package jwt

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	session "github.com/weisd/casbin-admin/middleware/jwt-session"
)

var (
	// ContextKey ContextKey
	ContextKey = "_JwtSessionKey"

	// DefaultConfig DefaultConfig
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

// Config Config
type Config struct {
	Skipper middleware.Skipper
	Manager *session.SessionManager
}

// Middleware Middleware
func Middleware(m *session.SessionManager) echo.MiddlewareFunc {
	config := DefaultConfig
	config.Manager = m

	return MiddlewareWithConfig(config)
}

// MiddlewareWithConfig MiddlewareWithConfig
func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return h(c)
			}

			sess, err := config.Manager.GetSession(c.Request())
			if err != nil {
				return err
			}

			cliams := sess.GetCliams()

			c.Logger().Debug("cliams:", cliams.Author, cliams.Data, cliams.Verify, cliams.CAA)

			// 验证

			c.Set(ContextKey, cliams)

			return h(c)

		}
	}

}

// Cliams Cliams
func Cliams(c echo.Context) *session.SessionClaims {
	f := c.Get(ContextKey)
	if f == nil {
		return session.NewSessionClaims()
	}

	return f.(*session.SessionClaims)
}
