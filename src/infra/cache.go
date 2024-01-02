package infra

import (
	"sync"
	"time"
)

type cacheItem struct {
	expireIn time.Time
	value    interface{}
}

type Cache struct {
	mutex sync.RWMutex
	items map[string]cacheItem
}

func NewCache() *Cache {
	return &Cache{items: map[string]cacheItem{}}
}

func (ref *Cache) get(key string) (item cacheItem, exists bool) {
	ref.mutex.RLock()
	item, exists = ref.items[key]
	ref.mutex.RUnlock()
	return item, exists
}

func (ref *Cache) Get(key string) interface{} {
	cacheItem, exists := ref.get(key)
	if !exists {
		return nil
	}
	if cacheItem.expireIn.Before(time.Now()) {
		ref.mutex.Lock()
		delete(ref.items, key)
		ref.mutex.Unlock()
	}
	return cacheItem.value
}

func (ref *Cache) Save(key string, value interface{}, ttl time.Duration) {
	ref.mutex.Lock()
	defer ref.mutex.Unlock()
	if ref.items == nil {
		ref.items = map[string]cacheItem{}
	}
	ref.items[key] = cacheItem{
		expireIn: time.Now().Add(ttl),
		value:    value,
	}
}
