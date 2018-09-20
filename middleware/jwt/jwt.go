package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
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

			sess, err := config.Manager.GetSession(c.Request())
			if err != nil {
				return err
			}

			cliams := sess.GetCliams()

			c.Logger().Debug("cliams:", c.RealIP(), cliams.Author, cliams.Data, cliams.Verify, cliams.CAA)

			// 永远传递session
			c.Set(ContextKey, sess)

			// 过滤验证
			if config.Skipper(c) {
				return h(c)
			}

			// 验证
			valid, err := sess.Valid()
			if err != nil {
				return err
			}

			if config.Manager.Options().VerifyIP.True() && !cliams.VerifyKey("ip", c.RealIP()) {

				return jwt.NewValidationError("ip verify failed", jwt.ValidationErrorMalformed)
			}

			if !valid {
				return jwt.NewValidationError("jwt valid failed", jwt.ValidationErrorMalformed)
			}

			return h(c)

		}
	}

}

// Session Session
func Session(c echo.Context) *session.Session {
	f := c.Get(ContextKey)
	if f == nil {
		panic("session context not found")
	}

	return f.(*session.Session)
}
