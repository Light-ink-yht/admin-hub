package email_svc

import (
	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/email_repo"
)

// EmailService 邮件服务接口
type EmailService interface {
	// 邮件配置相关方法
	CreateEmailConfig(config *email_domain.EmailConfig) error
	UpdateEmailConfig(config *email_domain.EmailConfig) error
	DeleteEmailConfig(id string) error
	GetEmailConfigByID(id string) (*email_domain.EmailConfig, error)
	GetAllEmailConfigs() ([]*email_domain.EmailConfig, error)
	UpdateEmailConfigStatus(id string, status string) error
	GetDefaultEmailConfig() (*email_domain.EmailConfig, error)
	SetDefaultEmailConfig(id string) error

	// 邮件模板相关方法
	CreateEmailTemplate(template *email_domain.EmailTemplate) error
	UpdateEmailTemplate(template *email_domain.EmailTemplate) error
	DeleteEmailTemplate(id string) error
	GetEmailTemplateByID(id string) (*email_domain.EmailTemplate, error)
	GetAllEmailTemplates() ([]*email_domain.EmailTemplate, error)
	GetEmailTemplateByType(templateType string) (*email_domain.EmailTemplate, error)
	SearchEmailTemplates(keyword string) ([]*email_domain.EmailTemplate, error)
	UpdateEmailTemplateStatus(id string, status string) error
}

// emailService 邮件服务实现
type emailService struct {
	emailConfigRepo   email_repo.EmailConfigRepository
	emailTemplateRepo email_repo.EmailTemplateRepository
}

// NewEmailService 创建邮件服务实例
func NewEmailService(
	emailConfigRepo email_repo.EmailConfigRepository,
	emailTemplateRepo email_repo.EmailTemplateRepository,
) EmailService {
	return &emailService{
		emailConfigRepo:   emailConfigRepo,
		emailTemplateRepo: emailTemplateRepo,
	}
}

// CreateEmailConfig 创建邮件配置
func (s *emailService) CreateEmailConfig(config *email_domain.EmailConfig) error {
	return s.emailConfigRepo.Create(config)
}

// UpdateEmailConfig 更新邮件配置
func (s *emailService) UpdateEmailConfig(config *email_domain.EmailConfig) error {
	return s.emailConfigRepo.Update(config)
}

// DeleteEmailConfig 删除邮件配置
func (s *emailService) DeleteEmailConfig(id string) error {
	return s.emailConfigRepo.Delete(id)
}

// GetEmailConfigByID 根据ID获取邮件配置
func (s *emailService) GetEmailConfigByID(id string) (*email_domain.EmailConfig, error) {
	return s.emailConfigRepo.GetByID(id)
}

// GetAllEmailConfigs 获取所有邮件配置
func (s *emailService) GetAllEmailConfigs() ([]*email_domain.EmailConfig, error) {
	return s.emailConfigRepo.GetAll()
}

// UpdateEmailConfigStatus 更新邮件配置状态
func (s *emailService) UpdateEmailConfigStatus(id string, status string) error {
	if status == "enable" {
		return s.emailConfigRepo.Enable(id)
	} else if status == "disable" {
		return s.emailConfigRepo.Disable(id)
	}
	return nil
}

// GetDefaultEmailConfig 获取默认邮件配置
func (s *emailService) GetDefaultEmailConfig() (*email_domain.EmailConfig, error) {
	return s.emailConfigRepo.GetDefault()
}

// SetDefaultEmailConfig 设置默认邮件配置
func (s *emailService) SetDefaultEmailConfig(id string) error {
	return s.emailConfigRepo.SetDefault(id)
}

// CreateEmailTemplate 创建邮件模板
func (s *emailService) CreateEmailTemplate(template *email_domain.EmailTemplate) error {
	return s.emailTemplateRepo.Create(template)
}

// UpdateEmailTemplate 更新邮件模板
func (s *emailService) UpdateEmailTemplate(template *email_domain.EmailTemplate) error {
	return s.emailTemplateRepo.Update(template)
}

// DeleteEmailTemplate 删除邮件模板
func (s *emailService) DeleteEmailTemplate(id string) error {
	return s.emailTemplateRepo.Delete(id)
}

// GetEmailTemplateByID 根据ID获取邮件模板
func (s *emailService) GetEmailTemplateByID(id string) (*email_domain.EmailTemplate, error) {
	return s.emailTemplateRepo.GetByID(id)
}

// GetAllEmailTemplates 获取所有邮件模板
func (s *emailService) GetAllEmailTemplates() ([]*email_domain.EmailTemplate, error) {
	return s.emailTemplateRepo.GetAll()
}

// GetEmailTemplateByType 根据类型获取邮件模板
func (s *emailService) GetEmailTemplateByType(templateType string) (*email_domain.EmailTemplate, error) {
	return s.emailTemplateRepo.GetByType(templateType)
}

// SearchEmailTemplates 搜索邮件模板
func (s *emailService) SearchEmailTemplates(keyword string) ([]*email_domain.EmailTemplate, error) {
	return s.emailTemplateRepo.Search(keyword)
}

// UpdateEmailTemplateStatus 更新邮件模板状态
func (s *emailService) UpdateEmailTemplateStatus(id string, status string) error {
	if status == "enable" {
		return s.emailTemplateRepo.Enable(id)
	} else if status == "disable" {
		return s.emailTemplateRepo.Disable(id)
	}
	return nil
}
