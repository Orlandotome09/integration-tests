package _interfaces

import "time"

type Cache interface {
	Get(key string) interface{}
	Save(key string, value interface{}, ttl time.Duration)
}
