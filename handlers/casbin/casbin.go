package casbin

import (
	"net/http"

	"github.com/labstack/echo"
	en "github.com/weisd/casbin-admin/models/casbin"
)

// Index Index
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

// PermissionList 权限列表 p
func PermissionList(c echo.Context) error {
	return c.JSON(http.StatusOK, en.Enfc.GetPolicy())
}

// PermissionAdd PermissionAdd
func PermissionAdd(c echo.Context) error {

	args := &en.Permission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("PermissionAdd", args)

	return c.JSON(http.StatusOK, en.Enfc.PermissionAdd(args))
}

// PermissionDel PermissionDel
func PermissionDel(c echo.Context) error {

	args := &en.Permission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("PermissionDel", args)

	return c.JSON(http.StatusOK, en.Enfc.PermissionDel(args))
}

// RoleList 角色列表 p
func RoleList(c echo.Context) error {

	return c.JSON(http.StatusOK, en.Enfc.GetRoles())
}

// RolePermissions g 角色对应权限列表
func RolePermissions(c echo.Context) error {
	role := c.FormValue("role")

	return c.JSON(http.StatusOK, en.Enfc.GetRolesForUser(role))
}

// RolePermissionAdd RolePermissionAdd
func RolePermissionAdd(c echo.Context) error {
	args := &en.RolePermission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("RolePermissionAdd", args)

	return c.JSON(http.StatusOK, en.Enfc.AddRoleForUser(args.Role, args.Permission))
}

// RolePermissionDel RolePermissionDel
func RolePermissionDel(c echo.Context) error {
	args := &en.RolePermission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("RolePermissionDel", args)

	return c.JSON(http.StatusOK, en.Enfc.DeleteRoleForUser(args.Role, args.Permission))
}

// UserRoles UserRoles
func UserRoles(c echo.Context) error {

	uid := c.FormValue("uid")

	return c.JSON(http.StatusOK, en.Enfc.GetRolesForUser(uid))
}

// UserRoleAdd UserRoleAdd
func UserRoleAdd(c echo.Context) error {
	args := &en.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("UserRoleAdd", args)

	return c.JSON(http.StatusOK, en.Enfc.AddRoleForUser(args.UID, args.Role))
}

// UserRoleDel UserRoleDel
func UserRoleDel(c echo.Context) error {
	args := &en.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("UserRoleDel", args)

	return c.JSON(http.StatusOK, en.Enfc.DeleteRoleForUser(args.UID, args.Role))
}

// UserPermissions UserPermissions
func UserPermissions(c echo.Context) error {

	uid := c.FormValue("uid")

	// uid->role
	roles := en.Enfc.GetRolesForUser(uid)

	permissions := make([]*en.Permission, 0)

	for i := range roles {

		//  role->permission
		rolePermissions := en.Enfc.GetRolesForUser(roles[i])

		for j := range rolePermissions {
			ps := en.Enfc.GetPermissionsForUser(rolePermissions[j])
			for k := range ps {
				p := &en.Permission{}
				p.Parse(ps[k])
				permissions = append(permissions, p)
			}

		}

	}

	return c.JSON(http.StatusOK, permissions)
}

// Enforce Enforce
func Enforce(c echo.Context) error {
	args := &en.Request{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("Enforce", args)

	return c.JSON(http.StatusOK, en.Enfc.Enforce(args))
}

// GetModel GetModel
func GetModel(c echo.Context) error {
	return c.JSON(http.StatusOK, en.Enfc.GetModel())
}
