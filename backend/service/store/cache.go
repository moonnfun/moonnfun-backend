package store

import (
	"errors"
	"fmt"
	"meme3/global"
	"sync"
	"time"
)

func NewCache() *sync.Map {
	return new(sync.Map)
}

func CachePush(vcache *sync.Map, key, val any, bForce bool) error {
	if vcache == nil {
		return errors.New("invalid cache")
	}
	if !bForce {
		if _, ok := vcache.Load(key); ok {
			return fmt.Errorf("cache data already exists, key: %v", key)
		}
	}
	vcache.Store(key, val)
	return nil
}

// delete: beforeDelleteFunc return true
func CachePop(vcache *sync.Map, key any, bDelete bool, beforeDelleteFunc func(k, v any) bool) any {
	if vcache == nil {
		return errors.New("invalid cache")
	}
	val, ok := vcache.Load(key)
	if ok && bDelete {
		if beforeDelleteFunc != nil && !beforeDelleteFunc(key, val) {
			return val
		}
		vcache.Delete(key)
	}
	return val
}

func CachePushByTime(vcache *sync.Map, key any, val any, bForce bool, timeout time.Duration, callback func(key, val any) bool) error {
	if vcache == nil {
		return errors.New("invalid cache")
	}
	if !bForce {
		if _, ok := vcache.Load(key); ok {
			return fmt.Errorf("cache data already exists, key: %v", key)
		}
	}
	vcache.Store(key, val)

	if timeout > 0 {
		go func(vcache *sync.Map, key, val any, timeout time.Duration, callback func(key, val any) bool) {
			for {
				select {
				case <-time.After(timeout):
					global.Debug("before remove memory data", "key", key, "timeout", timeout)
					CachePop(vcache, key, true, callback)
					return
				}
			}
		}(vcache, key, val, timeout, callback)
	}
	return nil
}
