package jwt

import (
	"sync"
)

// MemeryStore MemeryStore
type MemeryStore struct {
	lock *sync.RWMutex
	data map[int64]CaaItem
}

// NewMemeryStore NewMemeryStore
func NewMemeryStore() *MemeryStore {
	return &MemeryStore{
		lock: new(sync.RWMutex),
		data: make(map[int64]CaaItem),
	}
}

// Init Init
func (p *MemeryStore) Init(config string) error {
	return nil
}

// GetCounter GetCounter
func (p *MemeryStore) GetCounter(uid int64) (int64, error) {
	p.lock.RLock()
	c, has := p.data[uid]
	p.lock.RUnlock()
	if !has {
		return 0, nil
	}
	return c.Counter, nil
}

// SetCounter SetCounter
func (p *MemeryStore) SetCounter(uid int64, n int64) error {
	p.lock.Lock()
	c, has := p.data[uid]
	if !has {
		c = CaaItem{}
	}

	c.Counter = n
	p.data[uid] = c
	p.lock.Unlock()

	return nil
}

// GetTimeout GetTimeout
func (p *MemeryStore) GetTimeout(uid int64) (int64, error) {
	p.lock.RLock()
	c, has := p.data[uid]
	p.lock.RUnlock()
	if !has {
		return 0, nil
	}
	return c.Timeout, nil
}

// SetTimeout SetTimeout
func (p *MemeryStore) SetTimeout(uid int64, t int64) error {

	p.lock.Lock()
	c, has := p.data[uid]
	if !has {
		c = CaaItem{}
	}

	c.Timeout = t
	p.data[uid] = c
	p.lock.Unlock()

	return nil
}

func init() {
	RegisterStore("memery", NewMemeryStore())
}
