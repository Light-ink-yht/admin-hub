package code_svc

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/code_repo"
	"github.com/Light-ink-yht/admin-hub/internal/repository/sms_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/msg_svc/sms"
	"github.com/Light-ink-yht/admin-hub/internal/service/msg_svc/sms/tencent"
)

var codeTplId = "1877556"

type SmsServiceFace interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type SmsService struct {
	repo            *code_repo.CodeRepository
	logService      log_svc.LogService
	smsConfigRepo   sms_repo.SmsConfigRepository
	smsTemplateRepo sms_repo.SmsTemplateRepository
	smsSvc          sms.SmsService
}

func NewSmsService(
	repo *code_repo.CodeRepository,
	logService log_svc.LogService,
	smsConfigRepo sms_repo.SmsConfigRepository,
	smsTemplateRepo sms_repo.SmsTemplateRepository,
) SmsServiceFace {
	return &SmsService{
		repo:            repo,
		logService:      logService,
		smsConfigRepo:   smsConfigRepo,
		smsTemplateRepo: smsTemplateRepo,
	}
}

// Send 发验证码，我需要什么参数？
func (svc *SmsService) Send(ctx context.Context, biz string, phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis 和数据库
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 记录失败日志
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("发送验证码失败: %s", phone), "失败", err.Error(), "", nil)
		return err
	}

	// 使用默认模板ID
	codeTplId := codeTplId

	// 从数据库获取短信模板
	smsTemplate, err := svc.smsTemplateRepo.FindByIdentifier(ctx, "code_template")
	if err != nil {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("获取短信模板失败"), "失败", err.Error(), "", nil)
		return fmt.Errorf("获取短信模板失败: %w", err)
	}
	if smsTemplate != nil {
		// 使用数据库中的模板ID
		codeTplId = smsTemplate.SMB004
	} else {
		// 如果没有找到对应的模板，使用默认的硬编码模板作为备份
		svc.logService.LogBusiness(ctx, log_domain.LogLevelWarn, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("未找到短信模板: code_template, 使用默认模板"), "警告", "", "", nil)
	}

	// 从数据库获取短信服务商配置
	smsConfig, err := svc.smsConfigRepo.GetDefault(ctx)
	if err != nil {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("获取短信配置失败"), "失败", err.Error(), "", nil)
		return fmt.Errorf("获取短信配置失败: %w", err)
	}
	if smsConfig == nil || !smsConfig.IsEnabled() {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("未找到启用的短信服务商配置"), "失败", "", "", nil)
		return fmt.Errorf("未找到启用的短信服务商配置")
	}

	// 根据服务商类型创建对应的短信服务实例
	switch smsConfig.GetProviderType() {
	case "tencent":
		// 创建腾讯云短信服务实例
		tencentSmsService := tencent.NewService(nil, smsConfig.SMA004, smsConfig.SMA007)
		svc.smsSvc = tencentSmsService
	default:
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("不支持的短信服务商类型: %s", smsConfig.GetProviderType()), "失败", "", "", nil)
		return fmt.Errorf("不支持的短信服务商类型: %s", smsConfig.GetProviderType())
	}

	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)

	// 记录业务日志
	if err != nil {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelError, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("短信发送失败: %s", phone), "失败", err.Error(), "", nil)
	} else {
		svc.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "验证码服务", "发送验证码", "", "", "",
			fmt.Sprintf("短信发送成功: %s :%s", phone, code), "成功", "", "", nil)
	}

	return err
}

func (svc *SmsService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

// generateCode 生成一个六位数的验证码
func (svc *SmsService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}
