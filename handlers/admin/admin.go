package admin

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"github.com/weisd/casbin-admin/handlers"
	"github.com/weisd/casbin-admin/middleware/jwt"
	session "github.com/weisd/casbin-admin/middleware/jwt-session"
	"github.com/weisd/casbin-admin/models/admin"
)

// Login Login
func Login(c echo.Context) error {
	args := &AccountPasswd{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Login.Bind")
	}

	c.Logger().Debug("Login.args", args)

	info, err := admin.InfoByAccount(args.Account)
	if err != nil {
		return errors.Wrap(err, "Login.InfoByAccount")
	}
	if info == nil {
		return handlers.NewError(404, "User Not Found or Wrong Password")
	}

	c.Logger().Debug("Login.info", info)

	if !admin.CheckPasswd(info, args.Passwd) {
		return handlers.NewError(404, "User Not Found or Wrong Password1")
	}

	sess := jwt.Session(c)
	// sess := session.Manager.NewSession()

	sess.GetCliams().SetAuthor(session.Author{ID: info.ID, Name: info.Name})

	token, err := sess.SignedString()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, token)
}

// Logout Logout
func Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, "Logout")
}

// Add Add
func Add(c echo.Context) error {

	args := &admin.CasbinAdmin{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Add.Bind")
	}

	// @TODO verify empty value
	if err := handlers.ValidateStruct(args); err != nil {
		// errs := err.(validator.ValidationErrors)
		return handlers.NewError(403, err.Error())
	}

	if len(args.Phone) != 0 {
		old, err := admin.InfoByPhone(args.Phone)
		if err != nil {
			return errors.Wrap(err, "Add.InfoByPhone")
		}
		if old != nil {
			return handlers.NewError(403, "Phone exists")
		}
	}

	if len(args.Email) != 0 {
		old, err := admin.InfoByEmail(args.Email)
		if err != nil {
			return errors.Wrap(err, "Add.InfoByEmail")
		}
		if old != nil {
			return handlers.NewError(403, "Email exists")
		}
	}

	if len(args.Name) != 0 {
		old, err := admin.InfoByName(args.Name)
		if err != nil {
			return errors.Wrap(err, "Add.InfoByName")
		}
		if old != nil {
			return handlers.NewError(403, "Name exists")
		}
	}

	err := admin.Create(args)
	if err != nil {
		return errors.Wrap(err, "Add.Create")
	}

	return c.JSON(http.StatusOK, args)
}

// UpdateName UpdateName
func UpdateName(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateName")
}

// UpdatePwd UpdatePwd
func UpdatePwd(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdatePwd")
}
