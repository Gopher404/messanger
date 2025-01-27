package cache

import (
	"messanger/pkg/errors"
	"time"
)

type Cache struct {
	c map[string]value
}

type value struct {
	v int
	t time.Time
}

func NewCache() *Cache {
	c := &Cache{
		c: make(map[string]value),
	}

	go func() {
		for {
			time.Sleep(time.Minute * 5)
			now := time.Now()
			for k, v := range c.c {
				if now.After(v.t) {
					delete(c.c, k)
				}
			}
		}
	}()

	return c
}

func (c *Cache) Set(key string, v int, ttl time.Duration) *errors.Error {
	c.c[key] = value{
		v: v,
		t: time.Now().Add(ttl),
	}
	return nil
}

func (c *Cache) Get(key string) (int, *errors.Error) {
	v, ok := c.c[key]
	if ok {
		if time.Now().Before(v.t) {
			return v.v, nil
		} else {
			delete(c.c, key)
		}
	}

	return 0, nil
}

func (c *Cache) Del(key string) *errors.Error {
	delete(c.c, key)
	return nil
}
