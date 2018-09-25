package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	session "github.com/gocommon/jwt-session"
	"github.com/gocommon/jwt-session-middleware/echo"
	"github.com/labstack/echo"
	"github.com/weisd/casbin-admin/handlers"
	"github.com/weisd/casbin-admin/models"
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
		return handlers.NewError(401, "User Not Found or Wrong Password")
	}

	c.Logger().Debug("Login.info", info)

	if !admin.CheckPasswd(info, args.Passwd) {
		return handlers.NewError(401, "User Not Found or Wrong Password1")
	}

	sess := jwt.Session(c)

	sess.SetAuthor(session.Author{ID: info.ID, Name: info.Name, Avator: "https://avatars3.githubusercontent.com/u/2057561?s=460&v=4"})

	token, err := sess.SignedString()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResp{Token: token})
}

// Logout Logout
func Logout(c echo.Context) error {
	sess := jwt.Session(c)
	c.Logger().Info(sess.Author())
	sess.Flush()
	return c.JSON(http.StatusOK, nil)
}

// Info Info
func Info(c echo.Context) error {
	args := &ID{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Info.Bind")
	}

	info, err := admin.InfoByID(args.ID)
	if err != nil {
		return err
	}

	if err := handlers.ValidateStruct(args); err != nil {
		// errs := err.(validator.ValidationErrors)
		return handlers.NewError(403, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

// Add Add
func Add(c echo.Context) error {

	args := &CasbinAdminCreate{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Add.Bind")
	}

	c.Logger().Debug(args)

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

	info := &admin.CasbinAdmin{}
	info.Name = args.Name
	info.Email = args.Email
	info.Phone = args.Phone
	info.Status = args.Status
	info.Passwd = args.Passwd

	err := admin.Create(info)
	if err != nil {
		return errors.Wrap(err, "Add.Create")
	}

	return c.JSON(http.StatusOK, args)
}

// Update Update
func Update(c echo.Context) error {
	args := &CasbinAdminUpdate{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Update.Bind")
	}

	// @TODO verify empty value
	if err := handlers.ValidateStruct(args); err != nil {
		// errs := err.(validator.ValidationErrors)
		return handlers.NewError(403, err.Error())
	}

	info, err := admin.InfoByID(args.ID)
	if err != nil {
		return err
	} else if info == nil {
		return handlers.NewError(404, "info not found")
	}

	info.Name = args.Name
	info.Email = args.Email
	info.Phone = args.Phone
	info.Status = args.Status

	err = admin.Update(info)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info)
}

// UpdateName UpdateName
func UpdateName(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateName")
}

// UpdatePwd UpdatePwd
func UpdatePwd(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdatePwd")
}

// Index Index
func Index(c echo.Context) error {

	args := &CasbinAdminSearchArgs{}
	if err := c.Bind(args); err != nil {
		return errors.Wrap(err, "Update.Bind")
	}

	c.Logger().Info("CasbinAdminSearchArgs", args)

	var wheres []string
	var values []interface{}

	for i := range args.Querys {

		arr := strings.Split(args.Querys[i], ",")
		if len(arr) < 3 {
			return handlers.NewError(400, "bad requery querys")
		}

		// check fields
		if !models.ValidField(arr[0]) {
			return handlers.NewError(400, "bad requery querys1")
		}

		// check cond
		if !models.InConditions(arr[1]) {
			return handlers.NewError(400, "bad requery querys3")
		}

		condition, argsLen := models.MakeCondition(arr[1])
		wheres = append(wheres, fmt.Sprintf("%s %s", arr[0], condition))

		if len(arr[2:]) < argsLen {
			return handlers.NewError(400, "bad requery querys2")
		}

		switch arr[1] {
		case models.ConditionIN, models.ConditionNOTIN:
			values = append(values, arr[2:])
		default:
			for _, v := range arr[2 : 2+argsLen] {
				values = append(values, v)
			}

		}
	}

	// c.Logger().Info("=>", strings.Join(wheres, " and "), values)

	ex := models.SQLEx{}
	ex.Limit = args.Limit
	ex.Order = args.Order

	c.Logger().Info(ex)

	list, err := admin.ListSearch(strings.Join(wheres, " and "), values, ex)
	if err != nil {
		return err
	}

	for i := 0; i < len(list); i++ {
		c.Logger().Info(list[i])
	}

	return c.JSON(http.StatusOK, list)
}
