package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	ErrLocalCacheCodeSendTooMany   = errors.New("发送太频繁")
	ErrLocalCacheCodeVerifyTooMany = errors.New("验证太频繁")
)

type LocalCacheCodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type MemoryCodeCache struct {
	c *cache.Cache
}

func NewLocalCacheCode(defaultExpiration, cleanupInterval time.Duration) CodeCache {
	return &MemoryCodeCache{
		c: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *MemoryCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := c.key(biz, phone)
	if _, found := c.c.Get(key); found {
		return ErrCodeSendTooMany
	}

	// Set the cache with the provided code and a default expiration time
	c.c.Set(key, code, cache.DefaultExpiration)
	return nil
}

func (c *MemoryCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	key := c.key(biz, phone)
	value, found := c.c.Get(key)
	if !found {
		return false, nil
	}

	// Verify the code
	if value != code {
		return false, ErrCodeVerifyTooMany
	}

	// Optionally delete the code after verification
	c.c.Delete(key)
	return true, nil
}

func (c *MemoryCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
