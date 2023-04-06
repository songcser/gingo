package utils

import (
	"encoding/json"
	"github.com/coocood/freecache"
)

type Cache[T any] struct {
	cache *freecache.Cache
}

func NewCache[T any](size int) *Cache[T] {
	c := &Cache[T]{
		cache: freecache.NewCache(size),
	}
	return c
}

func (c *Cache[T]) Set(key string, val T) error {
	data, _ := json.Marshal(val)
	return c.cache.Set([]byte(key), data, 60)
}

func (c *Cache[T]) Get(key string) (T, error) {
	var t T
	v, err := c.cache.Get([]byte(key))
	if nil == err {
		err = json.Unmarshal(v, &t)
		if nil == err {
			return t, nil
		}
	}
	return t, err
}
