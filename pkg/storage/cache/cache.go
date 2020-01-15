package cache

import (
	"fmt"
	"strings"
	"sync"
)

// key is "api_usd_currency"
type stored map[string]Details

type Details struct {
	Price           string
	ChangePCTHour   string
	ChangePCT24Hour string
	ChangePCT7Day   string
}

type Cache struct {
	sync.RWMutex
	items stored
}

type Key struct {
	API      string
	Fiat     string
	Currency string
}

func NewCache() *Cache {
	return &Cache{
		items: make(stored),
	}
}

func (c *Cache) Set(k Key, d Details) {
	c.Lock()
	c.items[key(k)] = d
	c.Unlock()
}

func GenKey(a, f, c string) Key {
	return Key{API: strings.ToLower(a), Fiat: strings.ToLower(f), Currency: strings.ToLower(c)}
}

func key(k Key) string {
	return fmt.Sprintf("%s_%s_%s", k.API, k.Fiat, k.Currency)
}

func (c *Cache) Get(k Key) (d Details, ok bool) {
	c.RLock()
	defer c.RUnlock()
	d, ok = c.items[key(k)]
	if !ok {
		return Details{}, false
	}
	return
}

func (c *Cache) Delete(k Key) {
	c.Lock()
	delete(c.items, key(k))
	c.Unlock()
}
