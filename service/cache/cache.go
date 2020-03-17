package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	CacheDefaultExpiration    time.Duration = 5 * time.Minute
	SendCodeDefaultExpiration time.Duration = 1 * time.Minute
)

// 缓存数据
var (
	OC = NewCache()
)

/**
 * 新建缓存Cache
 */
func NewCache() *cache.Cache {
	oc := cache.New(5*time.Minute, 1*time.Minute)
	return oc
}
