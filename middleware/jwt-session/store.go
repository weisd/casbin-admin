package jwt

import (
	"encoding/binary"
	"fmt"
)

// CaaItem CaaItem
type CaaItem struct {
	Counter int64
	Timeout int64
}

// DecodeCaaItem DecodeCaaItem
func DecodeCaaItem(b []byte) CaaItem {
	if len(b) != 16 {
		return CaaItem{}
	}

	c := CaaItem{}

	c.Counter = int64(binary.LittleEndian.Uint64(b[0:]))
	c.Timeout = int64(binary.LittleEndian.Uint64(b[8:]))

	return c
}

// EncodeCaaItem EncodeCaaItem
func EncodeCaaItem(c CaaItem) []byte {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b[0:], uint64(c.Counter))
	binary.LittleEndian.PutUint64(b[8:], uint64(c.Timeout))
	return b
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
