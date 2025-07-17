package store

import (
	"fmt"
	"meme3/global"
	"sync"
	"time"
)

var vCache sync.Map
var DisableCache bool

func Cache() *sync.Map {
	return &vCache
}

func CacheListPageKey(key any, total int64, limit, offset int) string {
	return fmt.Sprintf("%v-%v-%v-%v", key, total, limit, offset)
}

func CacheSet(key, val any, bForce bool) error {
	if !bForce {
		if _, ok := vCache.Load(key); ok {
			return fmt.Errorf("cache data already exists, key: %v", key)
		}
	}
	vCache.Store(key, val)
	return nil
}

// delete: beforeDelleteFunc return true
func CacheGet(key any, bDelete bool, beforeDelleteFunc func(v any) bool) any {
	val, ok := vCache.Load(key)
	if ok && bDelete {
		if beforeDelleteFunc != nil && !beforeDelleteFunc(val) {
			return val
		}
		vCache.Delete(key)
	}
	return val
}

func CacheSetByTime(key any, val any, bForce bool, timeout time.Duration, callback func(val any) bool) error {
	if !bForce {
		if _, ok := vCache.Load(key); ok {
			return fmt.Errorf("cache data already exists, key: %v", key)
		}
	}
	vCache.Store(key, val)

	if timeout > 0 {
		go autoClearByTimer(key, val, timeout, callback)
	}
	return nil
}

// timeUnit: second
func autoClearByTimer(key, val any, timeout time.Duration, callback func(val any) bool) {
	for {
		select {
		case <-time.After(timeout):
			global.Debug("before remove memory data", "key", key, "timeout", timeout)
			CacheGet(key, true, callback)
			return
		}
	}
}
