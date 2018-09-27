package casbin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/weisd/casbin-admin/models/casbin"
)

// Index Index
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

// PermissionList 权限列表 p
func PermissionList(c echo.Context) error {
	arr := casbin.Enfc.GetPolicy()

	c.Logger().Info(arr)
	list := make([]*casbin.Permission, len(arr))
	for i := range arr {
		info := &casbin.Permission{}
		info.Parse(arr[i])
		list[i] = info

	}
	return c.JSON(http.StatusOK, list)
}

// PermissionAdd PermissionAdd
func PermissionAdd(c echo.Context) error {

	args := &casbin.Permission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("PermissionAdd", args)

	return c.JSON(http.StatusOK, casbin.Enfc.PermissionAdd(args))
}

// PermissionDel PermissionDel
func PermissionDel(c echo.Context) error {

	args := &casbin.Permission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("PermissionDel", args)

	//  同时删除角色关联
	existsRoles := casbin.Enfc.RolePermissionListByPermission(args.Name)
	for i := range existsRoles {
		casbin.Enfc.DeleteRoleForUser(existsRoles[i].Role, existsRoles[i].Permission)
	}

	return c.JSON(http.StatusOK, casbin.Enfc.PermissionDel(args))
}

// RoleList 角色列表 p 返回role名称列表
func RoleList(c echo.Context) error {

	list := []casbin.RolePermissionList{}

	roleList := casbin.Enfc.GetRoles()
	for _, v := range roleList {
		info := casbin.RolePermissionList{}
		info.Role = v
		info.PermissionList = casbin.Enfc.RolePermissionList(v)
		list = append(list, info)
	}

	return c.JSON(http.StatusOK, list)
}

// RolePermissions g 角色对应权限列表 返回permission列表
func RolePermissions(c echo.Context) error {
	args := &casbin.RolePermission{}
	if err := c.Bind(args); err != nil {
		return err
	}
	// GetNamedPolicy

	return c.JSON(http.StatusOK, casbin.Enfc.RolePermissionList(args.Role))
}

// RolePermissionAdd RolePermissionAdd
func RolePermissionAdd(c echo.Context) error {
	args := &casbin.RolePermission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	permissionList := strings.Split(args.Permission, ",")
	if len(permissionList) == 0 {
		return echo.NewHTTPError(400, "至少选择一个权限")
	}

	// 取旧值
	oldPermissionList := casbin.Enfc.RolePermissionList(args.Role)

	delList := []string{}

	for _, oldInfo := range oldPermissionList {
		var in bool
		for _, newName := range permissionList {
			if oldInfo.Name == newName {
				in = true
			}
		}

		if !in {
			delList = append(delList, oldInfo.Name)
		}
	}

	addList := []string{}

	for _, newName := range permissionList {
		var in bool
		for _, oldInfo := range oldPermissionList {
			if oldInfo.Name == newName {
				in = true
			}
		}

		if !in {
			addList = append(addList, newName)
		}
	}

	for i := range addList {
		if !casbin.Enfc.PermissionExists(permissionList[i]) {
			return echo.NewHTTPError(400, "paermission not exists:"+permissionList[i])
		}

		casbin.Enfc.AddRoleForUser(args.Role, addList[i])
	}

	for i := range delList {
		if !casbin.Enfc.PermissionExists(permissionList[i]) {
			return echo.NewHTTPError(400, "paermission not exists:"+permissionList[i])
		}

		casbin.Enfc.DeleteRoleForUser(args.Role, delList[i])
	}

	return c.JSON(http.StatusOK, casbin.Enfc.RolePermissionList(args.Role))
}

// RolePermissionDel RolePermissionDel
func RolePermissionDel(c echo.Context) error {
	args := &casbin.RolePermission{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("RolePermissionDel", args)

	return c.JSON(http.StatusOK, casbin.Enfc.DeleteRoleForUser(args.Role, args.Permission))
}

// RoleUserList RoleUserList 所有用户角色
func RoleUserList(c echo.Context) error {
	return c.JSON(http.StatusOK, casbin.Enfc.GetRoleUsers())
}

// UserRoles UserRoles
func UserRoles(c echo.Context) error {

	args := &casbin.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("UserRoles args", args)

	return c.JSON(http.StatusOK, casbin.Enfc.GetRolesForUser(strconv.FormatInt(args.UID, 10)))
}

// UserRoleAdd UserRoleAdd
func UserRoleAdd(c echo.Context) error {
	args := &casbin.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	roleList := strings.Split(args.Role, ",")

	sUID := strconv.FormatInt(args.UID, 10)

	// 取旧值
	oldList := casbin.Enfc.GetRolesForUser(sUID)

	delList := []string{}

	for _, oldInfo := range oldList {
		var in bool
		for _, newName := range roleList {
			if oldInfo == newName {
				in = true
			}
		}

		if !in {
			delList = append(delList, oldInfo)
		}
	}

	addList := []string{}

	for _, newName := range roleList {
		var in bool
		for _, oldInfo := range oldList {
			if oldInfo == newName {
				in = true
			}
		}

		if !in {
			addList = append(addList, newName)
		}
	}

	for i := range addList {
		if !casbin.Enfc.RoleExists(roleList[i]) {
			return echo.NewHTTPError(400, "paermission not exists:"+roleList[i])
		}

		casbin.Enfc.AddRoleForUser(sUID, addList[i])
	}

	for i := range delList {
		if !casbin.Enfc.RoleExists(roleList[i]) {
			return echo.NewHTTPError(400, "paermission not exists:"+roleList[i])
		}

		casbin.Enfc.DeleteRoleForUser(sUID, delList[i])
	}

	c.Logger().Debug("UserRoleAdd", args)

	return c.JSON(http.StatusOK, casbin.Enfc.GetRolesForUser(sUID))
}

// RoleDel RoleDel
func RoleDel(c echo.Context) error {
	args := &casbin.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	// 删除权限关联

	// 删除用户关联

	return c.JSON(http.StatusOK, casbin.Enfc.DeleteRoleByName(args.Role))
}

// UserPermissions UserPermissions
func UserPermissions(c echo.Context) error {

	args := &casbin.UserRole{}
	if err := c.Bind(args); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, casbin.Enfc.UserPermissions(args.UID))

	// uid->role
	// roles := casbin.Enfc.GetRolesForUser(uid)

	// permissions := make([]*casbin.Permission, 0)

	// for i := range roles {

	// 	//  role->permission
	// 	rolePermissions := casbin.Enfc.GetRolesForUser(roles[i])

	// 	for j := range rolePermissions {
	// 		ps := casbin.Enfc.GetPermissionsForUser(rolePermissions[j])
	// 		for k := range ps {
	// 			p := &casbin.Permission{}
	// 			p.Parse(ps[k])
	// 			permissions = append(permissions, p)
	// 		}

	// 	}

	// }

	// return c.JSON(http.StatusOK, permissions)
}

// Enforce Enforce
func Enforce(c echo.Context) error {
	args := &casbin.Request{}
	if err := c.Bind(args); err != nil {
		return err
	}

	c.Logger().Debug("Enforce", args)

	return c.JSON(http.StatusOK, casbin.Enfc.Enforce(args))
}

// GetModel GetModel
func GetModel(c echo.Context) error {
	return c.JSON(http.StatusOK, casbin.Enfc.GetModel())
}
