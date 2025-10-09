package email_repo

import (
	"context"
	"errors"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
	"github.com/google/uuid"
)

// EmailConfigRepositoryImpl 邮件服务器配置仓库实现
// 实现EmailConfigRepository接口，负责邮件服务器配置的业务逻辑处理
// 包含邮件服务器配置的增删改查、启用/禁用、默认配置设置等功能

type EmailConfigRepositoryImpl struct {
	emailConfigDAO *dao.EmailConfigDAO // 邮件服务器配置DAO对象
}

// NewEmailConfigRepository 创建新的邮件服务器配置仓库实例
// 参数：
// - emailConfigDAO: 邮件服务器配置DAO对象
// 返回值：
// - EmailConfigRepository: 邮件服务器配置仓库接口实现
func NewEmailConfigRepository(emailConfigDAO *dao.EmailConfigDAO) EmailConfigRepository {
	return &EmailConfigRepositoryImpl{
		emailConfigDAO: emailConfigDAO,
	}
}

// Create 创建邮件服务器配置
// 创建新的邮件服务器配置记录
// 参数：
// - config: 邮件服务器配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) Create(config *email_domain.EmailConfig) error {
	// 创建上下文
	ctx := context.Background()
	// 检查必要字段是否为空
	if config.EMA002 == "" || config.EMA003 == "" || config.EMA004 == "" || config.EMA006 == "" {
		return errors.New("发送方邮箱地址、授权码、SMTP服务器地址和配置名称不能为空")
	}

	// 生成主键ID
	if config.EMA001 == "" {
		config.EMA001 = uuid.New().String()
	}

	// 设置创建和更新时间
	timeNow := time.Now()
	if config.EMA008.IsZero() {
		config.EMA008 = timeNow
	}
	config.EMA009 = timeNow

	// 插入数据到数据库
	return repo.emailConfigDAO.Insert(ctx, config)
}

// Update 更新邮件服务器配置
// 更新邮件服务器配置记录
// 参数：
// - config: 邮件服务器配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) Update(config *email_domain.EmailConfig) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if config.EMA001 == "" {
		return errors.New("配置ID不能为空")
	}

	// 检查配置是否存在
	existingConfig, err := repo.emailConfigDAO.FindByID(ctx, config.EMA001)
	if err != nil {
		return err
	}
	if existingConfig == nil {
		return errors.New("配置不存在")
	}

	// 检查必要字段是否为空
	if config.EMA002 == "" || config.EMA003 == "" || config.EMA004 == "" || config.EMA006 == "" {
		return errors.New("发送方邮箱地址、授权码、SMTP服务器地址和配置名称不能为空")
	}

	// 设置更新时间
	config.EMA009 = time.Now()

	// 更新数据到数据库
	return repo.emailConfigDAO.Update(ctx, config)
}

// Delete 删除邮件服务器配置
// 删除邮件服务器配置记录
// 参数：
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) Delete(id string) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return errors.New("配置ID不能为空")
	}

	// 检查配置是否存在
	existingConfig, err := repo.emailConfigDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingConfig == nil {
		return errors.New("配置不存在")
	}

	// 检查是否为默认配置
	if existingConfig.EMA010 == "default" {
		return errors.New("默认配置不能删除")
	}

	// 从数据库删除
	return repo.emailConfigDAO.Delete(ctx, id)
}

// GetByID 根据ID获取邮件服务器配置
// 获取指定ID的邮件服务器配置信息
// 参数：
// - id: 邮件服务器配置ID
// 返回值：
// - *email_domain.EmailConfig: 邮件服务器配置领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) GetByID(id string) (*email_domain.EmailConfig, error) {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return nil, errors.New("配置ID不能为空")
	}

	// 从数据库查询
	return repo.emailConfigDAO.FindByID(ctx, id)
}

// GetAll 获取所有邮件服务器配置
// 获取所有邮件服务器配置信息
// 返回值：
// - []*email_domain.EmailConfig: 邮件服务器配置领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) GetAll() ([]*email_domain.EmailConfig, error) {
	// 创建上下文
	ctx := context.Background()
	// 从数据库查询所有配置
	return repo.emailConfigDAO.FindAll(ctx)
}

// Enable 启用邮件服务器配置
// 启用指定ID的邮件服务器配置
// 参数：
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) Enable(id string) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return errors.New("配置ID不能为空")
	}

	// 检查配置是否存在
	existingConfig, err := repo.emailConfigDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingConfig == nil {
		return errors.New("配置不存在")
	}

	// 启用配置
	return repo.emailConfigDAO.Enable(ctx, id)
}

// Disable 禁用邮件服务器配置
// 禁用指定ID的邮件服务器配置
// 参数：
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) Disable(id string) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return errors.New("配置ID不能为空")
	}

	// 检查配置是否存在
	existingConfig, err := repo.emailConfigDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingConfig == nil {
		return errors.New("配置不存在")
	}

	// 检查是否为默认配置
	if existingConfig.EMA010 == "default" {
		return errors.New("默认配置不能禁用")
	}

	// 禁用配置
	return repo.emailConfigDAO.Disable(ctx, id)
}

// GetDefault 获取默认邮件服务器配置
// 获取系统默认的邮件服务器配置
// 返回值：
// - *email_domain.EmailConfig: 邮件服务器配置领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) GetDefault() (*email_domain.EmailConfig, error) {
	// 创建上下文
	ctx := context.Background()
	// 从数据库查询默认配置
	return repo.emailConfigDAO.FindDefault(ctx)
}

// SetDefault 设置默认邮件服务器配置
// 设置指定ID的邮件服务器配置为默认配置
// 参数：
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailConfigRepositoryImpl) SetDefault(id string) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return errors.New("配置ID不能为空")
	}

	// 检查配置是否存在
	existingConfig, err := repo.emailConfigDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingConfig == nil {
		return errors.New("配置不存在")
	}

	// 检查配置是否已启用
	if existingConfig.EMA007 != 1 {
		return errors.New("只能设置启用状态的配置为默认配置")
	}

	// 设置为默认配置
	return repo.emailConfigDAO.SetDefault(ctx, id)
}

// EmailTemplateRepositoryImpl 邮件模板仓库实现
// 实现EmailTemplateRepository接口，负责邮件模板的业务逻辑处理
// 包含邮件模板的增删改查、启用/禁用等功能

type EmailTemplateRepositoryImpl struct {
	emailTemplateDAO *dao.EmailTemplateDAO // 邮件模板DAO对象
}

// NewEmailTemplateRepository 创建新的邮件模板仓库实例
// 参数：
// - emailTemplateDAO: 邮件模板DAO对象
// 返回值：
// - EmailTemplateRepository: 邮件模板仓库接口实现
func NewEmailTemplateRepository(emailTemplateDAO *dao.EmailTemplateDAO) EmailTemplateRepository {
	return &EmailTemplateRepositoryImpl{
		emailTemplateDAO: emailTemplateDAO,
	}
}

// Create 创建邮件模板
// 创建新的邮件模板记录
// 参数：
// - template: 邮件模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Create(template *email_domain.EmailTemplate) error {
	// 创建上下文
	ctx := context.Background()
	// 检查必要字段是否为空
	if template.EMB002 == "" || template.EMB003 == "" || template.EMB004 == "" || template.EMB005 == "" {
		return errors.New("模板名称、模板类型、邮件主题和邮件内容不能为空")
	}

	// 生成主键ID
	if template.EMB001 == "" {
		template.EMB001 = uuid.New().String()
	}

	// 设置创建和更新时间
	timeNow := time.Now()
	if template.EMB007.IsZero() {
		template.EMB007 = timeNow
	}
	template.EMB008 = timeNow

	// 插入数据到数据库
	return repo.emailTemplateDAO.Insert(ctx, template)
}

// Update 更新邮件模板
// 更新邮件模板记录
// 参数：
// - template: 邮件模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Update(template *email_domain.EmailTemplate) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if template.EMB001 == "" {
		return errors.New("模板ID不能为空")
	}

	// 检查模板是否存在
	existingTemplate, err := repo.emailTemplateDAO.FindByID(ctx, template.EMB001)
	if err != nil {
		return err
	}
	if existingTemplate == nil {
		return errors.New("模板不存在")
	}

	// 检查必要字段是否为空
	if template.EMB002 == "" || template.EMB003 == "" || template.EMB004 == "" || template.EMB005 == "" {
		return errors.New("模板名称、模板类型、邮件主题和邮件内容不能为空")
	}

	// 设置更新时间
	template.EMB008 = time.Now()

	// 更新数据到数据库
	return repo.emailTemplateDAO.Update(ctx, template)
}

// Delete 删除邮件模板
// 删除邮件模板记录
// 参数：
// - id: 邮件模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Delete(id string) error {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return errors.New("模板ID不能为空")
	}

	// 检查模板是否存在
	existingTemplate, err := repo.emailTemplateDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTemplate == nil {
		return errors.New("模板不存在")
	}

	// 从数据库删除
	return repo.emailTemplateDAO.Delete(ctx, id)
}

// GetByID 根据ID获取邮件模板
// 获取指定ID的邮件模板信息
// 参数：
// - id: 邮件模板ID
// 返回值：
// - *email_domain.EmailTemplate: 邮件模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) GetByID(id string) (*email_domain.EmailTemplate, error) {
	// 创建上下文
	ctx := context.Background()
	// 检查ID是否为空
	if id == "" {
		return nil, errors.New("模板ID不能为空")
	}

	// 从数据库查询
	return repo.emailTemplateDAO.FindByID(ctx, id)
}

// GetAll 获取所有邮件模板
// 获取所有邮件模板信息
// 返回值：
// - []*email_domain.EmailTemplate: 邮件模板领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) GetAll() ([]*email_domain.EmailTemplate, error) {
	// 创建上下文
	ctx := context.Background()
	// 从数据库查询所有模板
	return repo.emailTemplateDAO.FindAll(ctx)
}

// GetByType 根据类型获取邮件模板
// 获取指定类型的邮件模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - templateType: 模板标识（EMB003字段）
// 返回值：
// - *email_domain.EmailTemplate: 邮件模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) GetByType(templateType string) (*email_domain.EmailTemplate, error) {
	// 检查类型是否为空
	if templateType == "" {
		return nil, errors.New("模板标识不能为空")
	}

	// 创建上下文
	ctx := context.Background()
	// 从数据库查询指定标识的模板
	return repo.emailTemplateDAO.FindByType(ctx, templateType)
}

// Enable 启用邮件模板
// 启用指定ID的邮件模板
// 参数：
// - id: 邮件模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Enable(id string) error {
	// 检查ID是否为空
	if id == "" {
		return errors.New("模板ID不能为空")
	}

	// 创建上下文
	ctx := context.Background()
	// 检查模板是否存在
	template, err := repo.emailTemplateDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return errors.New("模板不存在")
	}

	// 启用模板
	template.EMB006 = 1
	template.EMB008 = time.Now()
	return repo.emailTemplateDAO.Update(ctx, template)
}

// Disable 禁用邮件模板
// 禁用指定ID的邮件模板
// 参数：
// - id: 邮件模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Disable(id string) error {
	// 检查ID是否为空
	if id == "" {
		return errors.New("模板ID不能为空")
	}

	// 创建上下文
	ctx := context.Background()
	// 检查模板是否存在
	template, err := repo.emailTemplateDAO.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return errors.New("模板不存在")
	}

	// 禁用模板
	template.EMB006 = 0
	template.EMB008 = time.Now()
	return repo.emailTemplateDAO.Update(ctx, template)
}

// Search 搜索邮件模板
// 根据关键词搜索邮件模板
// 参数：
// - keyword: 搜索关键词
// 返回值：
// - []*email_domain.EmailTemplate: 匹配的邮件模板领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (repo *EmailTemplateRepositoryImpl) Search(keyword string) ([]*email_domain.EmailTemplate, error) {
	// 创建上下文
	ctx := context.Background()
	// 从数据库搜索模板
	return repo.emailTemplateDAO.Search(ctx, keyword)
}
