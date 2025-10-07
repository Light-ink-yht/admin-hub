package sms_repo

import (
	"context"

	"github.com/Light-ink-yht/admin-hub/internal/domain/sms_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
)

// smsConfigRepository 短信配置仓库实现
type smsConfigRepository struct {
	dao *dao.SmsConfigDAO
}

// NewSmsConfigRepository 创建短信配置仓库实例
func NewSmsConfigRepository(dao *dao.SmsConfigDAO) SmsConfigRepository {
	return &smsConfigRepository{
		dao: dao,
	}
}

// Save 保存短信配置
func (r *smsConfigRepository) Save(ctx context.Context, config *sms_domain.SmsConfig) error {
	return r.dao.InsertConfig(ctx, config)
}

// Update 更新短信配置
func (r *smsConfigRepository) Update(ctx context.Context, config *sms_domain.SmsConfig) error {
	return r.dao.UpdateConfig(ctx, config)
}

// Delete 删除短信配置
func (r *smsConfigRepository) Delete(ctx context.Context, id string) error {
	return r.dao.DeleteConfig(ctx, id)
}

// FindByID 根据ID查找短信配置
func (r *smsConfigRepository) FindByID(ctx context.Context, id string) (*sms_domain.SmsConfig, error) {
	return r.dao.FindConfigByID(ctx, id)
}

// FindAll 查询所有短信配置
func (r *smsConfigRepository) FindAll(ctx context.Context) ([]*sms_domain.SmsConfig, error) {
	return r.dao.FindAllConfigs(ctx)
}

// FindEnabled 查询所有启用的短信配置
func (r *smsConfigRepository) FindEnabled(ctx context.Context) ([]*sms_domain.SmsConfig, error) {
	return r.dao.FindEnabledConfigs(ctx)
}

// Enable 启用短信配置
func (r *smsConfigRepository) Enable(ctx context.Context, id string) error {
	return r.dao.EnableConfig(ctx, id)
}

// Disable 禁用短信配置
func (r *smsConfigRepository) Disable(ctx context.Context, id string) error {
	return r.dao.DisableConfig(ctx, id)
}

// GetDefault 获取默认短信配置（第一个启用的配置）
func (r *smsConfigRepository) GetDefault(ctx context.Context) (*sms_domain.SmsConfig, error) {
	configs, err := r.dao.FindEnabledConfigs(ctx)
	if err != nil {
		return nil, err
	}

	// 返回第一个启用的配置作为默认配置
	if len(configs) > 0 {
		return configs[0], nil
	}

	return nil, nil
}

// smsTemplateRepository 短信模板仓库实现
type smsTemplateRepository struct {
	dao *dao.SmsConfigDAO
}

// NewSmsTemplateRepository 创建短信模板仓库实例
func NewSmsTemplateRepository(dao *dao.SmsConfigDAO) SmsTemplateRepository {
	return &smsTemplateRepository{
		dao: dao,
	}
}

// Save 保存短信模板
func (r *smsTemplateRepository) Save(ctx context.Context, template *sms_domain.SmsTemplate) error {
	return r.dao.InsertTemplate(ctx, template)
}

// Update 更新短信模板
func (r *smsTemplateRepository) Update(ctx context.Context, template *sms_domain.SmsTemplate) error {
	return r.dao.UpdateTemplate(ctx, template)
}

// Delete 删除短信模板
func (r *smsTemplateRepository) Delete(ctx context.Context, id string) error {
	return r.dao.DeleteTemplate(ctx, id)
}

// FindByID 根据ID查找短信模板
func (r *smsTemplateRepository) FindByID(ctx context.Context, id string) (*sms_domain.SmsTemplate, error) {
	return r.dao.FindTemplateByID(ctx, id)
}

// FindByIdentifier 根据标识查找短信模板
func (r *smsTemplateRepository) FindByIdentifier(ctx context.Context, identifier string) (*sms_domain.SmsTemplate, error) {
	return r.dao.FindTemplateByIdentifier(ctx, identifier)
}

// FindAll 查询所有短信模板
func (r *smsTemplateRepository) FindAll(ctx context.Context) ([]*sms_domain.SmsTemplate, error) {
	return r.dao.FindAllTemplates(ctx)
}

// FindEnabled 查询所有启用的短信模板
func (r *smsTemplateRepository) FindEnabled(ctx context.Context) ([]*sms_domain.SmsTemplate, error) {
	return r.dao.FindEnabledTemplates(ctx)
}

// Enable 启用短信模板
func (r *smsTemplateRepository) Enable(ctx context.Context, id string) error {
	return r.dao.EnableTemplate(ctx, id)
}

// Disable 禁用短信模板
func (r *smsTemplateRepository) Disable(ctx context.Context, id string) error {
	return r.dao.DisableTemplate(ctx, id)
}
