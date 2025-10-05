package code_svc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/code_repo"
	"github.com/Light-ink-yht/admin-hub/internal/repository/email_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/msg_svc/email/tencent"
)

type EmailServiceFace interface {
	Send(ctx context.Context, biz string, email string, template string) error
	Verify(ctx context.Context, biz string, email string, inputCode string) (bool, error)
}

type EmailService struct {
	repo              *code_repo.CodeRepository
	logService        log_svc.LogService
	emailConfigRepo   email_repo.EmailConfigRepository
	emailTemplateRepo email_repo.EmailTemplateRepository
}

func NewEmailService(
	repo *code_repo.CodeRepository,
	logService log_svc.LogService,
	emailConfigRepo email_repo.EmailConfigRepository,
	emailTemplateRepo email_repo.EmailTemplateRepository,
) EmailServiceFace {
	return &EmailService{
		repo:              repo,
		logService:        logService,
		emailConfigRepo:   emailConfigRepo,
		emailTemplateRepo: emailTemplateRepo,
	}
}

// Send 发送验证码 biz 区分业务场景
func (svc *EmailService) Send(ctx context.Context, biz string, email string, templateType string) error {
	// 生成验证码
	code := svc.generateCode()

	// 塞进 redis 和数据库
	err := svc.repo.Store(ctx, biz, email, code)
	if err != nil {
		// 记录失败日志
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("发送验证码失败: %s", email), "失败", err.Error(), "", nil)
		return err
	}

	// 从数据库获取邮件模板
	emailTemplate, err := svc.emailTemplateRepo.GetByType(templateType)
	if err != nil {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("获取邮件模板失败: %s", templateType), "失败", err.Error(), "", nil)
		return fmt.Errorf("获取邮件模板失败: %w", err)
	}
	if emailTemplate == nil {
		// 如果没有找到对应的模板，使用默认的硬编码模板作为备份
		svc.logService.LogBusiness(ctx, log_domain.LogLevelWarn, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("未找到邮件模板: %s, 使用默认模板", templateType), "警告", "", "", nil)
		emailTemplate = &email_domain.EmailTemplate{
			EMB004: "【验证码】您的验证码已生成",
			EMB005: "<p>您的验证码是：<strong>{code}</strong></p><p>请在10分钟内使用该验证码完成验证。</p>",
			EMB009: "默认验证码模板",
		}
	}

	// 替换模板中的 {code} 占位符
	body := strings.ReplaceAll(emailTemplate.EMB005, "{code}", code)

	// 从数据库获取邮件服务器配置
	emailConfig, err := svc.emailConfigRepo.GetDefault()
	if err != nil {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("获取邮件配置失败"), "失败", err.Error(), "", nil)
		return fmt.Errorf("获取邮件配置失败: %w", err)
	}
	if emailConfig == nil || !emailConfig.IsEnabled() {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("未找到启用的邮件服务器配置"), "失败", "", "", nil)
		return errors.New("未找到启用的邮件服务器配置")
	}

	// 使用数据库配置创建邮件服务实例
	emailService := tencent.NewEmailService(emailConfig.EMA002, emailConfig.EMA003, emailConfig.EMA004, emailConfig.EMA005)
	// 发送出去
	err = emailService.Send(ctx, email, emailTemplate.EMB004, body)

	// 记录业务日志
	if err != nil {
		err := svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("邮件发送失败: %s", email), "失败", err.Error(), "", nil)
		if err != nil {
			return err
		}
	} else {
		err := svc.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("邮件发送成功: %s :%s", email, code), "成功", "", "", nil)
		if err != nil {
			return err
		}
	}

	return err
}

// Verify 验证验证码
func (svc *EmailService) Verify(ctx context.Context, biz string, email string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, email, inputCode)
}

// generateCode 生成一个六位数的验证码
func (svc *EmailService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}
