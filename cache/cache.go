package cache

import (
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var defaultExpire = 1 * time.Hour
var defaultCache = NewCache(defaultExpire, 30*time.Minute)

func NewCache(defaultExpiration, cleanupInterval time.Duration) *gocache.Cache {
	return gocache.New(defaultExpiration, cleanupInterval)
}

func Add(key string, v interface{}, expire time.Duration) error {
	return defaultCache.Add(key, v, expire)
}

func Set(key string, v interface{}, expire time.Duration) {
	defaultCache.Set(key, v, expire)
}

func Get(key string) (interface{}, bool) {
	return defaultCache.Get(key)
}
