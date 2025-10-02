package code_repo

import (
	"context"
	"github.com/Light-ink-yht/admin-hub/internal/repository/cache/code_cache"
)

var ErrCodeSendTooMany = code_cache.ErrCodeSendTooMany
var RrrCodeVerifyTooMany = code_cache.RrrCodeVerifyTooMany

type CodeRepository struct {
	cache *code_cache.CodeCache
}

func NewCodeRepository(cache *code_cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

// Store 存到 Store 中，biz 区分业务场景，input 用户的phone/email，code 验证码
func (repo *CodeRepository) Store(ctx context.Context, biz, input, code string) error {
	return repo.cache.Set(ctx, biz, input, code)
}

// Verify Redis 中验证，biz 区分业务场景，input 用户的phone/email，code 验证码
func (repo *CodeRepository) Verify(ctx context.Context, biz, input, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, input, code)
}
