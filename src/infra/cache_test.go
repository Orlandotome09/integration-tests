package infra

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func Test_Get_Should_Return_Item_And_Delete_If_After(t *testing.T) {

	cache := Cache{
		mutex: sync.RWMutex{},
		items: map[string]cacheItem{
			"SomeKey": {
				expireIn: time.Now().Add(time.Minute * -1),
				value:    "SomeValue",
			},
		},
	}

	content := cache.Get("SomeKey")

	assert.Equal(t, "SomeValue", content)
	assert.Zero(t, len(cache.items))
}

func Test_Get_Should_Return_Existing_Item(t *testing.T) {

	cache := Cache{
		mutex: sync.RWMutex{},
		items: map[string]cacheItem{
			"SomeKey": {
				expireIn: time.Now().Add(time.Minute + 1),
				value:    "SomeValue",
			},
		},
	}

	content := cache.Get("SomeKey")

	assert.Equal(t, "SomeValue", content)
	assert.Equal(t, 1, len(cache.items))
}

func Test_Get_Should_Return_Non_Existing_Item(t *testing.T) {

	cache := Cache{
		mutex: sync.RWMutex{},
		items: map[string]cacheItem{
			"SomeKey": {
				expireIn: time.Now().Add(time.Minute + 1),
				value:    "SomeValue",
			},
		},
	}

	content := cache.Get("SomeNonExistingKey")

	assert.Nil(t, content)
}

func Test_Save_Item(t *testing.T) {

	cache := Cache{
		mutex: sync.RWMutex{},
		items: map[string]cacheItem{
			"SomeKey": {
				expireIn: time.Now().Add(time.Minute + 1),
				value:    "SomeValue",
			},
		},
	}

	cache.Save("SomeOtherKey", "SomeValue", time.Minute)

	content, exists := cache.items["SomeOtherKey"]

	assert.Equal(t, "SomeValue", content.value)
	assert.True(t, exists)
}

func Test_Override_Item(t *testing.T) {

	cache := Cache{
		mutex: sync.RWMutex{},
		items: map[string]cacheItem{
			"SomeKey": {
				expireIn: time.Now().Add(time.Minute + 1),
				value:    "SomeValue",
			},
		},
	}

	cache.Save("SomeKey", "SomeOtherValue", time.Minute)

	content, exists := cache.items["SomeKey"]

	assert.Equal(t, "SomeOtherValue", content.value)
	assert.True(t, exists)
}
