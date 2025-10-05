package code_svc

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/code_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/msg_svc/email/tencent"
)

type EmailServiceFace interface {
	Send(ctx context.Context, biz string, email string, template string) error
	Verify(ctx context.Context, biz string, email string, inputCode string) (bool, error)
}

type EmailService struct {
	repo       *code_repo.CodeRepository
	logService log_svc.LogService
}

func NewEmailService(repo *code_repo.CodeRepository, logService log_svc.LogService) EmailServiceFace {
	return &EmailService{
		repo:       repo,
		logService: logService,
	}
}

// Send 发送验证码 biz 区分业务场景
func (svc *EmailService) Send(ctx context.Context, biz string, email string, template string) error {

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

	// 替换模板中的 {code} 占位符
	body := strings.ReplaceAll(template, "{code}", code)
	// 假设这些是你的邮箱认证信息和SMTP服务器设置
	authEmail := "1480224563@qq.com" // 发送方邮箱地址
	authPwd := "kglwbdxsfosrhhhi"    // 发送方邮箱授权码
	smtpHost := "smtp.qq.com"        // SMTP服务器地址
	smtpPort := 465                  // SMTP服务器端口

	// 使用NewService创建Service实例
	emailService := tencent.NewEmailService(authEmail, authPwd, smtpHost, smtpPort)
	// 发送出去
	err = emailService.Send(ctx, email, "【灵脑科技】", body)

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
