package sms_svc

import (
	"context"
	"fmt"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/sms_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/sms_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
)

// Service 短信服务
type Service struct {
	smsConfigRepo   sms_repo.SmsConfigRepository
	smsTemplateRepo sms_repo.SmsTemplateRepository
	logService      log_svc.LogService
}

// NewSmsService 创建短信服务实例
func NewSmsService(
	smsConfigRepo sms_repo.SmsConfigRepository,
	smsTemplateRepo sms_repo.SmsTemplateRepository,
	logService log_svc.LogService,
) *Service {
	return &Service{
		smsConfigRepo:   smsConfigRepo,
		smsTemplateRepo: smsTemplateRepo,
		logService:      logService,
	}
}

// CreateSmsConfig 创建短信配置
func (s *Service) CreateSmsConfig(ctx context.Context, config *sms_domain.SmsConfig) error {
	// 记录创建配置的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "创建短信配置", "", "", "",
		fmt.Sprintf("开始创建短信配置: %s", config.SMA002), "处理中", "", "", nil)

	// 保存配置到数据库
	err := s.smsConfigRepo.Save(ctx, config)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "创建短信配置", "", "", "",
			fmt.Sprintf("创建短信配置失败: %s", config.SMA002), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "创建短信配置", "", "", "",
		fmt.Sprintf("成功创建短信配置: %s", config.SMA002), "成功", "", "", nil)

	return nil
}

// UpdateSmsConfig 更新短信配置
func (s *Service) UpdateSmsConfig(ctx context.Context, config *sms_domain.SmsConfig) error {
	// 记录更新配置的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "更新短信配置", "", "", "",
		fmt.Sprintf("开始更新短信配置: %s", config.SMA002), "处理中", "", "", nil)

	// 更新配置到数据库
	err := s.smsConfigRepo.Update(ctx, config)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "更新短信配置", "", "", "",
			fmt.Sprintf("更新短信配置失败: %s", config.SMA002), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "更新短信配置", "", "", "",
		fmt.Sprintf("成功更新短信配置: %s", config.SMA002), "成功", "", "", nil)

	return nil
}

// DeleteSmsConfig 删除短信配置
func (s *Service) DeleteSmsConfig(ctx context.Context, id string) error {
	// 记录删除配置的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "删除短信配置", "", "", "",
		fmt.Sprintf("开始删除短信配置: %s", id), "处理中", "", "", nil)

	// 从数据库删除配置
	err := s.smsConfigRepo.Delete(ctx, id)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "删除短信配置", "", "", "",
			fmt.Sprintf("删除短信配置失败: %s", id), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "删除短信配置", "", "", "",
		fmt.Sprintf("成功删除短信配置: %s", id), "成功", "", "", nil)

	return nil
}

// GetSmsConfig 获取短信配置
func (s *Service) GetSmsConfig(ctx context.Context, id string) (*sms_domain.SmsConfig, error) {
	// 记录获取配置的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取短信配置", "", "", "",
		fmt.Sprintf("开始获取短信配置: %s", id), "处理中", "", "", nil)

	// 从数据库获取配置
	config, err := s.smsConfigRepo.FindByID(ctx, id)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "获取短信配置", "", "", "",
			fmt.Sprintf("获取短信配置失败: %s", id), "失败", err.Error(), "", nil)
		return nil, err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取短信配置", "", "", "",
		fmt.Sprintf("成功获取短信配置: %s", id), "成功", "", "", nil)

	return config, nil
}

// GetAllSmsConfigs 获取所有短信配置
func (s *Service) GetAllSmsConfigs(ctx context.Context) ([]*sms_domain.SmsConfig, error) {
	// 记录获取所有配置的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取所有短信配置", "", "", "",
		"开始获取所有短信配置", "处理中", "", "", nil)

	// 从数据库获取所有配置
	configs, err := s.smsConfigRepo.FindAll(ctx)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "获取所有短信配置", "", "", "",
			"获取所有短信配置失败", "失败", err.Error(), "", nil)
		return nil, err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取所有短信配置", "", "", "",
		fmt.Sprintf("成功获取所有短信配置，共%d条", len(configs)), "成功", "", "", nil)

	return configs, nil
}

// CreateSmsTemplate 创建短信模板
func (s *Service) CreateSmsTemplate(ctx context.Context, template *sms_domain.SmsTemplate) error {
	// 记录创建模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "创建短信模板", "", "", "",
		fmt.Sprintf("开始创建短信模板: %s", template.SMB002), "处理中", "", "", nil)

	// 保存模板到数据库
	err := s.smsTemplateRepo.Save(ctx, template)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "创建短信模板", "", "", "",
			fmt.Sprintf("创建短信模板失败: %s", template.SMB002), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "创建短信模板", "", "", "",
		fmt.Sprintf("成功创建短信模板: %s", template.SMB002), "成功", "", "", nil)

	return nil
}

// UpdateSmsTemplate 更新短信模板
func (s *Service) UpdateSmsTemplate(ctx context.Context, template *sms_domain.SmsTemplate) error {
	// 记录更新模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "更新短信模板", "", "", "",
		fmt.Sprintf("开始更新短信模板: %s", template.SMB002), "处理中", "", "", nil)

	// 更新模板到数据库
	err := s.smsTemplateRepo.Update(ctx, template)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "更新短信模板", "", "", "",
			fmt.Sprintf("更新短信模板失败: %s", template.SMB002), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "更新短信模板", "", "", "",
		fmt.Sprintf("成功更新短信模板: %s", template.SMB002), "成功", "", "", nil)

	return nil
}

// DeleteSmsTemplate 删除短信模板
func (s *Service) DeleteSmsTemplate(ctx context.Context, id string) error {
	// 记录删除模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "删除短信模板", "", "", "",
		fmt.Sprintf("开始删除短信模板: %s", id), "处理中", "", "", nil)

	// 从数据库删除模板
	err := s.smsTemplateRepo.Delete(ctx, id)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "删除短信模板", "", "", "",
			fmt.Sprintf("删除短信模板失败: %s", id), "失败", err.Error(), "", nil)
		return err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "删除短信模板", "", "", "",
		fmt.Sprintf("成功删除短信模板: %s", id), "成功", "", "", nil)

	return nil
}

// GetSmsTemplate 获取短信模板
func (s *Service) GetSmsTemplate(ctx context.Context, id string) (*sms_domain.SmsTemplate, error) {
	// 记录获取模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取短信模板", "", "", "",
		fmt.Sprintf("开始获取短信模板: %s", id), "处理中", "", "", nil)

	// 从数据库获取模板
	template, err := s.smsTemplateRepo.FindByID(ctx, id)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "获取短信模板", "", "", "",
			fmt.Sprintf("获取短信模板失败: %s", id), "失败", err.Error(), "", nil)
		return nil, err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取短信模板", "", "", "",
		fmt.Sprintf("成功获取短信模板: %s", id), "成功", "", "", nil)

	return template, nil
}

// GetSmsTemplateByIdentifier 根据标识符获取短信模板
func (s *Service) GetSmsTemplateByIdentifier(ctx context.Context, identifier string) (*sms_domain.SmsTemplate, error) {
	// 记录获取模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "根据标识符获取短信模板", "", "", "",
		fmt.Sprintf("开始根据标识符获取短信模板: %s", identifier), "处理中", "", "", nil)

	// 从数据库获取模板
	template, err := s.smsTemplateRepo.FindByIdentifier(ctx, identifier)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "根据标识符获取短信模板", "", "", "",
			fmt.Sprintf("根据标识符获取短信模板失败: %s", identifier), "失败", err.Error(), "", nil)
		return nil, err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "根据标识符获取短信模板", "", "", "",
		fmt.Sprintf("成功根据标识符获取短信模板: %s", identifier), "成功", "", "", nil)

	return template, nil
}

// GetAllSmsTemplates 获取所有短信模板
func (s *Service) GetAllSmsTemplates(ctx context.Context) ([]*sms_domain.SmsTemplate, error) {
	// 记录获取所有模板的日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取所有短信模板", "", "", "",
		"开始获取所有短信模板", "处理中", "", "", nil)

	// 从数据库获取所有模板
	templates, err := s.smsTemplateRepo.FindAll(ctx)
	if err != nil {
		// 记录失败日志
		s.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信服务", "获取所有短信模板", "", "", "",
			"获取所有短信模板失败", "失败", err.Error(), "", nil)
		return nil, err
	}

	// 记录成功日志
	s.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信服务", "获取所有短信模板", "", "", "",
		fmt.Sprintf("成功获取所有短信模板，共%d条", len(templates)), "成功", "", "", nil)

	return templates, nil
}
