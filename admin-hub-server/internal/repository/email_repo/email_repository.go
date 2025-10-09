package email_repo

import (
	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
)

// EmailConfigRepository 邮件服务器配置仓库接口
type EmailConfigRepository interface {
	// Create 创建邮件服务器配置
	Create(config *email_domain.EmailConfig) error
	// Update 更新邮件服务器配置
	Update(config *email_domain.EmailConfig) error
	// Delete 删除邮件服务器配置
	Delete(id string) error
	// GetByID 根据ID获取邮件服务器配置
	GetByID(id string) (*email_domain.EmailConfig, error)
	// GetAll 获取所有邮件服务器配置
	GetAll() ([]*email_domain.EmailConfig, error)
	// Enable 启用邮件服务器配置
	Enable(id string) error
	// Disable 禁用邮件服务器配置
	Disable(id string) error
	// GetDefault 获取默认邮件服务器配置
	GetDefault() (*email_domain.EmailConfig, error)
	// SetDefault 设置默认邮件服务器配置
	SetDefault(id string) error
}

// EmailTemplateRepository 邮件模板仓库接口
type EmailTemplateRepository interface {
	// Create 创建邮件模板
	Create(template *email_domain.EmailTemplate) error
	// Update 更新邮件模板
	Update(template *email_domain.EmailTemplate) error
	// Delete 删除邮件模板
	Delete(id string) error
	// GetByID 根据ID获取邮件模板
	GetByID(id string) (*email_domain.EmailTemplate, error)
	// GetAll 获取所有邮件模板
	GetAll() ([]*email_domain.EmailTemplate, error)
	// GetByType 根据模板标识获取邮件模板
	GetByType(templateType string) (*email_domain.EmailTemplate, error)
	// Search 搜索邮件模板
	Search(keyword string) ([]*email_domain.EmailTemplate, error)
	// Enable 启用邮件模板
	Enable(id string) error
	// Disable 禁用邮件模板
	Disable(id string) error
}
