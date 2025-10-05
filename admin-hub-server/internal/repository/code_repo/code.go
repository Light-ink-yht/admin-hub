package code_repo

import (
	"context"
	"errors"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/code_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/cache/code_cache"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
	"github.com/google/uuid"
)

var ErrCodeSendTooMany = code_cache.ErrCodeSendTooMany
var RrrCodeVerifyTooMany = code_cache.RrrCodeVerifyTooMany
var ErrCodeDBOperationFailed = errors.New("验证码数据库操作失败")

// CodeRepository 验证码仓库
// 同时支持Redis缓存和数据库持久化
// 实现了验证码的存储、验证等核心功能
// 在保持原有Redis高性能的同时，增加了数据库持久化能力

// CodeExpireTime 验证码有效期，与Redis中的设置保持一致（600秒）
const CodeExpireTime = 600 * time.Second

// MaxVerifyAttempts 验证码最大验证次数
const MaxVerifyAttempts = 3

// SendFrequencyWindow 发送频率限制时间窗口（秒）
const SendFrequencyWindow = 60 * time.Second

// MaxSendCountInWindow 发送频率限制次数
const MaxSendCountInWindow = 5

// CodeRepository 验证码仓库结构体
// 包含Redis缓存和数据库访问对象

type CodeRepository struct {
	cache   *code_cache.CodeCache
	codeDAO *dao.CodeDAO
}

// NewCodeRepository 创建新的验证码仓库实例
// 同时初始化Redis缓存和数据库访问对象
func NewCodeRepository(cache *code_cache.CodeCache, codeDAO *dao.CodeDAO) *CodeRepository {
	return &CodeRepository{
		cache:   cache,
		codeDAO: codeDAO,
	}
}

// Store 存到 Store 中，biz 区分业务场景，input 用户的phone/email，code 验证码
// 同时存储到Redis和数据库中
// 1. 先检查发送频率
// 2. 存储到Redis
// 3. 存储到数据库
// 任何一步失败都返回错误
func (repo *CodeRepository) Store(ctx context.Context, biz, input, code string) error {
	return repo.StoreWithRemark(ctx, biz, input, code, "")
}

// StoreWithRemark 存储验证码并添加备注信息
// 扩展了Store方法，支持添加备注信息
func (repo *CodeRepository) StoreWithRemark(ctx context.Context, biz, input, code string, remark string) error {
	// 检查发送频率（数据库层面）
	startTime := time.Now().Add(-SendFrequencyWindow)
	count, err := repo.codeDAO.CountByBizAndInputInTimeRange(ctx, biz, input, startTime)
	if err != nil {
		return ErrCodeDBOperationFailed
	}
	if count >= MaxSendCountInWindow {
		return ErrCodeSendTooMany
	}

	// 先存储到Redis
	if err := repo.cache.Set(ctx, biz, input, code); err != nil {
		return err
	}

	// 再存储到数据库
	codeModel := code_domain.NewCode(
		biz,
		input,
		code,
		time.Now().Add(CodeExpireTime),
		MaxVerifyAttempts,
	)
	// 生成唯一ID
	codeModel.CDA001 = uuid.New().String()
	// 设置备注信息
	codeModel.CDA010 = remark

	if err := repo.codeDAO.Insert(ctx, codeModel); err != nil {
		// 数据库存储失败时，尝试删除Redis中的记录以保持一致性
		// 注意：这是一个尽力而为的操作，不影响主要业务流程
		repo.cache.Delete(ctx, biz, input)
		return ErrCodeDBOperationFailed
	}

	return nil
}

// Verify Redis 中验证，biz 区分业务场景，input 用户的phone/email，code 验证码
// 同时验证Redis和数据库中的状态
// 1. 先在Redis中验证
// 2. 根据验证结果更新数据库状态
// 3. 返回最终验证结果
func (repo *CodeRepository) Verify(ctx context.Context, biz, input, code string) (bool, error) {
	// 先在Redis中验证
	result, err := repo.cache.Verify(ctx, biz, input, code)
	if err != nil {
		return false, err
	}

	// 如果验证成功，更新数据库状态
	if result {
		// 查询数据库中的验证码记录
		codeModel, err := repo.codeDAO.FindByBizAndInput(ctx, biz, input)
		if err != nil {
			return false, ErrCodeDBOperationFailed
		}

		// 如果找到了记录，更新状态
		if codeModel != nil && codeModel.CanVerify() {
			// 标记为验证成功
			codeModel.VerifySuccess()
			// 更新数据库
			err = repo.codeDAO.Update(ctx, codeModel)
			if err != nil {
				// 数据库更新失败不影响验证结果，但记录错误
			}
		}
	} else if err == nil {
		// 验证失败但没有错误（可能是验证码错误或已过期）
		// 查询数据库中的验证码记录
		codeModel, err := repo.codeDAO.FindByBizAndInput(ctx, biz, input)
		if err != nil {
			return false, ErrCodeDBOperationFailed
		}

		// 如果找到了记录，更新状态
		if codeModel != nil && codeModel.CanVerify() {
			// 标记为验证失败
			codeModel.VerifyFailed()
			// 更新数据库
			err = repo.codeDAO.Update(ctx, codeModel)
			if err != nil {
				// 数据库更新失败不影响验证结果，但记录错误
			}
		}
	}

	return result, err
}

// Delete 删除指定的验证码记录
// 同时从Redis和数据库中删除
func (repo *CodeRepository) Delete(ctx context.Context, biz, input string) error {
	// 从Redis中删除
	if err := repo.cache.Delete(ctx, biz, input); err != nil {
		return err
	}

	// 从数据库中更新状态为已失效
	codeModel, err := repo.codeDAO.FindByBizAndInput(ctx, biz, input)
	if err != nil {
		return ErrCodeDBOperationFailed
	}

	if codeModel != nil {
		codeModel.CDA009 = 4 // 已失效
		if err := repo.codeDAO.Update(ctx, codeModel); err != nil {
			return ErrCodeDBOperationFailed
		}
	}

	return nil
}
