package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/gommon/log"

	jwtpkg "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/weisd/casbin-admin/handlers/admin"
	_ "github.com/weisd/casbin-admin/handlers/casbin"
	"github.com/weisd/casbin-admin/middleware/jwt"
	session "github.com/weisd/casbin-admin/middleware/jwt-session"
	"github.com/weisd/casbin-admin/router"
)

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// e.Pre(middleware.AddTrailingSlash())

	session.Init(session.Options{
		JwtSigningMethod: jwtpkg.SigningMethodHS256,
		JwtPrivateKey:[]byte("123"),
	})

	e.Use(jwt.Middleware(session.Manager))

	e.Debug = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = he.Message
		} else if e.Debug {
			msg = err.Error()
		} else {
			msg = http.StatusText(code)
		}
		if _, ok := msg.(string); ok {
			msg = echo.Map{"message": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == echo.HEAD { // Issue #608
				if err := c.NoContent(code); err != nil {
					goto ERROR
				}
			} else {
				if err := c.JSON(code, msg); err != nil {
					goto ERROR
				}
			}
		}
	ERROR:
		e.Logger.Error(err)
	}

	router.Init(e)

	start()

	// e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("local.com")

	// Start server
	go func() {
		// if err := e.StartAutoTLS(":443"); err != nil {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
