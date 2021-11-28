package redis

import "github.com/go-redis/redis/v8"

// NewRedisRing - создаёт новое подключение к redis.
func NewRedisRing(addresses map[string]string) *redis.Ring {
	return redis.NewRing(&redis.RingOptions{
		Addrs: addresses,
	})
}
