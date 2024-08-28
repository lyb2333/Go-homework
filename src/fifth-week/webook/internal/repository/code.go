package repository

import (
	"context"
	"homework/src/fifth-week/webook/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
var ErrCodeSendTooMany = cache.ErrCodeSendTooMany

type FailedRequest struct {
	Biz   string
	Phone string
	Code  string
}

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
	DumpFailedRequest(ctx context.Context, biz, phone, code string) error
	GetFailedRequests(ctx context.Context) ([]FailedRequest, error)
}

type CachedCodeRepository struct {
	cache cache.InMemoryCodeCache
}

func NewCodeRepository(c *cache.InMemoryCodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (c *CachedCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}

func (c *CachedCodeRepository) DumpFailedRequest(ctx context.Context, biz, phone, code string) error {
	return c.cache.DumpFailedRequest(ctx, biz, phone, code)
}

func (c *CachedCodeRepository) GetFailedRequests(ctx context.Context) ([]FailedRequest, error) {
	failedRequests, err := c.cache.GetFailedRequests(ctx)
	if err != nil {
		return nil, err
	}
	return convertToFailedRequestSlice(failedRequests), nil
}

func convertToFailedRequestSlice(cacheFailedRequests []cache.FailedRequest) []FailedRequest {
	failedRequests := make([]FailedRequest, len(cacheFailedRequests))
	for i, req := range cacheFailedRequests {
		failedRequests[i] = FailedRequest{
			Biz:   req.Biz,
			Phone: req.Phone,
			Code:  req.Code,
		}
	}
	return failedRequests
}
