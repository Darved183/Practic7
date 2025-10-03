package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheItem
	mux  sync.RWMutex
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.data[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.data, key)
}

func (c *Cache) CleanExpired() {
	c.mux.Lock()
	defer c.mux.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiresAt) {
			delete(c.data, key)
		}
	}
}

func main() {
	cache := NewCache()

	go func() {
		cache.CleanExpired()
	}()

	cache.Set("1", "Пригожин Женя", 5*time.Second)
	cache.Set("2", "Михаил Ждунов", 10*time.Second)
	cache.Set("3", "Игорь", 15*time.Second)

	if value, found := cache.Get("1"); found {
		fmt.Printf("1 = %v\n", value)
	}

	fmt.Println("Ожидаем 6 секунд...")
	time.Sleep(6 * time.Second)

	if value, found := cache.Get("1"); found {
		fmt.Printf("1 = %v\n", value)
	} else {
		fmt.Println("1 не найден")
	}

	if value, found := cache.Get("2"); found {
		fmt.Printf("2 = %v\n", value)
	}

	cache.Delete("2")
	if _, found := cache.Get("2"); !found {
		fmt.Println("2 удален")
	}

}
