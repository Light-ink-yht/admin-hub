package sms_repo

import (
	"context"

	"github.com/Light-ink-yht/admin-hub/internal/domain/sms_domain"
)

// SmsConfigRepository 短信配置仓库接口
// 定义了短信配置相关的仓储操作方法
type SmsConfigRepository interface {
	// Save 保存短信配置
	Save(ctx context.Context, config *sms_domain.SmsConfig) error

	// Update 更新短信配置
	Update(ctx context.Context, config *sms_domain.SmsConfig) error

	// Delete 删除短信配置
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找短信配置
	FindByID(ctx context.Context, id string) (*sms_domain.SmsConfig, error)

	// FindAll 查询所有短信配置
	FindAll(ctx context.Context) ([]*sms_domain.SmsConfig, error)

	// FindEnabled 查询所有启用的短信配置
	FindEnabled(ctx context.Context) ([]*sms_domain.SmsConfig, error)

	// Enable 启用短信配置
	Enable(ctx context.Context, id string) error

	// Disable 禁用短信配置
	Disable(ctx context.Context, id string) error

	// GetDefault 获取默认短信配置（第一个启用的配置）
	GetDefault(ctx context.Context) (*sms_domain.SmsConfig, error)
}

// SmsTemplateRepository 短信模板仓库接口
// 定义了短信模板相关的仓储操作方法
type SmsTemplateRepository interface {
	// Save 保存短信模板
	Save(ctx context.Context, template *sms_domain.SmsTemplate) error

	// Update 更新短信模板
	Update(ctx context.Context, template *sms_domain.SmsTemplate) error

	// Delete 删除短信模板
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找短信模板
	FindByID(ctx context.Context, id string) (*sms_domain.SmsTemplate, error)

	// FindByIdentifier 根据标识查找短信模板
	FindByIdentifier(ctx context.Context, identifier string) (*sms_domain.SmsTemplate, error)

	// FindAll 查询所有短信模板
	FindAll(ctx context.Context) ([]*sms_domain.SmsTemplate, error)

	// FindEnabled 查询所有启用的短信模板
	FindEnabled(ctx context.Context) ([]*sms_domain.SmsTemplate, error)

	// Enable 启用短信模板
	Enable(ctx context.Context, id string) error

	// Disable 禁用短信模板
	Disable(ctx context.Context, id string) error
}
