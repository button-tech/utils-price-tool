package storage

import (
	"errors"
	"sync"
)

type Api string

type CryptoCurrency string

type Fiat string

type FiatMap map[Fiat]map[CryptoCurrency]*Details

type Stored map[Api]FiatMap

type Details struct {
	Price           string
	ChangePCTHour   string
	ChangePCT24Hour string
	ChangePCT7Day   string
}

type Cache struct {
	sync.RWMutex
	items Stored
}

func NewCache() *Cache {
	return &Cache{
		items: make(Stored),
	}
}

func (c *Cache) Set(a Api, f FiatMap) {
	c.Lock()

	if _, ok := c.items[a]; !ok {
		c.items[a] = map[Fiat]map[CryptoCurrency]*Details{}
	}

	for k, v := range f {
		c.items[a][k] = v
	}

	c.Unlock()
}

func (c *Cache) Get(a Api) (f FiatMap, err error) {
	c.RLock()
	f = c.items[a]
	if f == nil {
		err = errors.New("cache: nil")
	}
	c.RUnlock()
	return
}
