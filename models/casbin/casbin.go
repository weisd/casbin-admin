package casbin

import (
	"unicode"

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
func (e *Enforcer) GetRoles() []*RolePermission {

	list := make([]*RolePermission, 0)

	groupPolicys := e.SyncedEnforcer.GetNamedGroupingPolicy("g")
	for i := range groupPolicys {
		r := &RolePermission{}
		r.Parse(groupPolicys[i])

		if isInt(r.Role) {
			continue
		}

		list = append(list, r)
	}

	return list
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
