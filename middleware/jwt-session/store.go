package jwt

import (
	"fmt"
)

// CaaItem CaaItem
type CaaItem struct {
	Counter int64
	Timeout int64
}

// Store Store
type Store interface {
	Init(config string) error

	GetCounter(uid int64) (int64, error)
	SetCounter(uid int64, n int64) error

	GetTimeout(uid int64) (int64, error)
	SetTimeout(uid int64, t int64) error
}

var stores = map[string]Store{}

// RegisterStore RegisterStore
func RegisterStore(name string, s Store) {
	if _, has := stores[name]; has {
		panic(fmt.Sprintf("jwt session store exists %s", name))
	}

	stores[name] = s
}

// NewStore NewStore
func NewStore(name, config string) Store {
	s, has := stores[name]
	if !has {
		panic(fmt.Sprintf("store not exists ! forgot to register?"))
	}

	err := s.Init(config)
	if err != nil {
		panic(fmt.Sprintf("store %s Init faild : %v", name, err))
	}

	return s
}
