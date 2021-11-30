package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
)

type serviceCache struct {
	cache *cache.Cache
	ttl   time.Duration
}

// NewServiceCache - создает новый инстанс кэша сервисов в redis.
func NewServiceCache(redis *redis.Ring, maxSize uint64, ttl time.Duration) *serviceCache {
	return &serviceCache{
		cache: cache.New(&cache.Options{
			Redis:      redis,
			LocalCache: cache.NewTinyLFU(int(maxSize), ttl),
		}),
		ttl: ttl,
	}
}

func (c *serviceCache) Set(ctx context.Context, service *subscription.Service) error {
	return c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   c.getKey(service.ID),
		Value: *service,
		TTL:   c.ttl,
	})
}

func (c *serviceCache) Unset(ctx context.Context, serviceID uint64) error {
	return c.cache.Delete(ctx, c.getKey(serviceID))
}

func (c *serviceCache) Get(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	var service subscription.Service

	if err := c.cache.Get(ctx, c.getKey(serviceID), &service); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, nil
		}

		return nil, err
	}

	return &service, nil
}

func (c *serviceCache) getKey(serviceID uint64) string {
	return fmt.Sprintf("service:%d", serviceID)
}
