package casbin

import (
	"strconv"

	"github.com/casbin/casbin"
)

// Enforcer Enforcer
type Enforcer struct {
	*casbin.SyncedEnforcer
	Options Options
}

// NewEnforcer NewEnforcer
func NewEnforcer(opt ...Option) *Enforcer {

	opts := NewOptions()
	for _, o := range opt {
		o(&opts)
	}

	e := &Enforcer{}
	e.Options = opts

	if e.Options.adapter != nil {
		e.SyncedEnforcer = casbin.NewSyncedEnforcer(casbin.NewModel(e.Options.conf), e.Options.adapter)
	} else {
		e.SyncedEnforcer = casbin.NewSyncedEnforcer(casbin.NewModel(e.Options.conf))
	}

	if e.Options.watcher != nil {
		e.SyncedEnforcer.SetWatcher(e.Options.watcher)
	}

	if e.Options.syncDuration > 0 {
		e.SyncedEnforcer.StartAutoLoadPolicy(e.Options.syncDuration)
	}

	return e

}

// Enfc Enfc
var Enfc *Enforcer

// Init Init
func Init(opt ...Option) {
	Enfc = NewEnforcer(opt...)
}

// PermissionAdd PermissionAdd
func (e *Enforcer) PermissionAdd(p *Permission) bool {
	return e.AddPolicy(p.Name, p.Path, p.Method, p.Origin)
}

// PermissionDel PermissionDel
func (e *Enforcer) PermissionDel(p *Permission) bool {
	return e.RemovePolicy(p.Name, p.Path, p.Method, p.Origin)
}

// Enforce Enforce
func (e *Enforcer) Enforce(req *Request) bool {
	return e.SyncedEnforcer.Enforce(req.UID, req.Path, req.Method, req.Origin)
}

// GetRoles GetRoles
func (e *Enforcer) GetRoles() []string {

	m := map[string]struct{}{}

	groupPolicys := e.SyncedEnforcer.GetNamedGroupingPolicy("g")
	for i := range groupPolicys {
		r := &RolePermission{}
		r.Parse(groupPolicys[i])

		if isInt(r.Role) {
			continue
		}

		m[r.Role] = struct{}{}

	}

	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}

	return list
}

// GetRoleUsers GetRoleUsers
func (e *Enforcer) GetRoleUsers() map[int64][]string {

	userRolsList := make([]*RolePermission, 0)

	groupPolicys := e.SyncedEnforcer.GetNamedGroupingPolicy("g")
	for i := range groupPolicys {
		r := &RolePermission{}
		r.Parse(groupPolicys[i])

		if isInt(r.Role) {
			userRolsList = append(userRolsList, r)
		}
	}

	uidRoles := make(map[int64][]string)
	for i := range userRolsList {
		info := userRolsList[i]
		uid, err := strconv.ParseInt(info.Role, 10, 64)
		if err != nil {
			// 不是uid过滤掉不要
			continue
		}

		if _, has := uidRoles[uid]; !has {
			uidRoles[uid] = make([]string, 0)
		}

		uidRoles[uid] = append(uidRoles[uid], info.Permission)
	}

	return uidRoles
}

// RolePermissionList RolePermissionList
func (e *Enforcer) RolePermissionList(role string) []*Permission {

	list := make([]*Permission, 0)

	roles := e.GetRolesForUser(role)
	for _, v := range roles {
		permission := e.PermissionInfoByName(v)
		if permission != nil {
			list = append(list, permission)
		}
	}

	return list

}

// PermissionExists PermissionExists
func (e *Enforcer) PermissionExists(name string) bool {
	return len(e.SyncedEnforcer.GetFilteredPolicy(0, name)) > 0
}

// PermissionInfoByName PermissionInfoByName
func (e *Enforcer) PermissionInfoByName(name string) *Permission {
	permission := e.SyncedEnforcer.GetFilteredPolicy(0, name)

	if len(permission) > 0 {
		info := &Permission{}
		info.Parse(permission[0])
		return info
	}
	return nil
}

// RolePermissionListByPermission RolePermissionListByPermission
func (e *Enforcer) RolePermissionListByPermission(name string) []*RolePermission {
	permission := e.SyncedEnforcer.GetFilteredGroupingPolicy(1, name)

	list := make([]*RolePermission, len(permission))
	for i := range permission {
		info := &RolePermission{}
		info.Parse(permission[i])
		list[i] = info
	}
	return list
}

// RolePermissionListByRole RolePermissionListByRole
func (e *Enforcer) RolePermissionListByRole(name string) []*RolePermission {
	permission := e.SyncedEnforcer.GetFilteredGroupingPolicy(0, name)

	list := make([]*RolePermission, len(permission))
	for i := range permission {
		info := &RolePermission{}
		info.Parse(permission[i])
		list[i] = info
	}
	return list
}

// UIDRoleListByRole UIDRoleListByRole
func (e *Enforcer) UIDRoleListByRole(name string) []*UserRole {
	result := e.SyncedEnforcer.GetFilteredGroupingPolicy(1, name)

	list := make([]*UserRole, len(result))
	for i := range result {
		info := &UserRole{}
		info.Parse(result[i])
		list[i] = info
	}
	return list
}

// RoleExists RoleExists
func (e *Enforcer) RoleExists(rols string) bool {
	return len(e.SyncedEnforcer.GetFilteredGroupingPolicy(0, rols)) > 0
}

// DeleteRoleByName DeleteRoleByName
func (e *Enforcer) DeleteRoleByName(name string) bool {
	list := e.SyncedEnforcer.GetFilteredGroupingPolicy(0, name)
	for i := range list {
		e.SyncedEnforcer.DeleteRoleForUser(list[i][0], list[i][1])
	}

	list = e.SyncedEnforcer.GetFilteredGroupingPolicy(1, name)
	for i := range list {
		e.SyncedEnforcer.DeleteRoleForUser(list[i][0], list[i][1])
	}

	return true
}

// UserPermissions UserPermissions
func (e *Enforcer) UserPermissions(uid int64) []*Permission {
	permission := e.SyncedEnforcer.GetPermissionsForUser(strconv.FormatInt(uid, 10))
	list := make([]*Permission, len(permission))
	for i := range permission {
		info := &Permission{}
		info.Parse(permission[i])
		list[i] = info
	}
	return list
}
