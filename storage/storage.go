package storage

import (
	"sync"
)

type Api string

type CryptoCurrency string

type Fiat string

type Stored map[Api]map[Fiat]map[CryptoCurrency]*Details

type Details struct {
	Price           string
	ChangePCTHour   string
	ChangePCT24Hour string
}

type Cache struct {
	sync.RWMutex
	items Stored
}

type Storage interface {
	Set(a Api, cr map[Fiat]map[CryptoCurrency]*Details)
	Get() Stored
}

func NewCache() Storage {
	items := make(Stored)

	cache := Cache{
		items: items,
	}

	return &cache
}

//todo: complete
func (c *Cache) Set(a Api, cr map[Fiat]map[CryptoCurrency]*Details) {
	c.Lock()




	for k, v := range cr {
		if val, ok := c.items[a]; ok {
			val[k] = v
		} else {
			val[k] = v
		}

	}

	c.Unlock()
}

func (c *Cache) Get() Stored {
	c.RLock()
	defer c.RUnlock()
	return c.items
}
