package casbin

import (
	"time"

	"github.com/casbin/casbin/persist"
)

var _conf = `
[request_definition]
r = sub, obj, act, org

[policy_definition]
p = sub, obj, act, org

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && KeyMatch(r.act, p.act) && r.org == p.org
`

// Options Options
type Options struct {
	syncDuration time.Duration
	watcher      persist.Watcher
	adapter      persist.Adapter
	conf         string
}

// Option Option
type Option func(*Options)

// NewOptions NewOptions
func NewOptions() Options {
	return Options{
		conf: _conf,
		// SyncDuration: time.Second * 10,
	}
}

// WithConf WithConf
func WithConf(conf string) Option {
	return func(o *Options) {
		o.conf = conf
	}
}

// WithSyncDuration WithSyncDuration
func WithSyncDuration(d time.Duration) Option {
	return func(o *Options) {
		o.syncDuration = d
	}
}

// WithWatcher WithWatcher
func WithWatcher(watcher persist.Watcher) Option {
	return func(o *Options) {
		o.watcher = watcher
	}
}

// WithAdapter WithAdapter
func WithAdapter(adapter persist.Adapter) Option {
	return func(o *Options) {
		o.adapter = adapter
	}
}
