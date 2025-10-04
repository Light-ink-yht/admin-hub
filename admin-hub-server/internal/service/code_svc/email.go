package code_svc

import (
	"context"
	"fmt"
	"github.com/Light-ink-yht/admin-hub/internal/repository/code_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/msg_svc/email/tencent"
	"go.uber.org/zap"
	"math/rand"
	"strings"
)

type EmailServiceFace interface {
	Send(ctx context.Context, biz string, email string, template string) error
	Verify(ctx context.Context, biz string, email string, inputCode string) (bool, error)
}

type EmailService struct {
	repo   *code_repo.CodeRepository
	logger *zap.Logger
}

func NewEmailService(repo *code_repo.CodeRepository, logger *zap.Logger) EmailServiceFace {
	return &EmailService{
		repo:   repo,
		logger: logger,
	}
}

// Send 发送验证码 biz 区分业务场景
func (svc *EmailService) Send(ctx context.Context, biz string, email string, template string) error {
	svc.logger.Info("开始发送验证码", zap.String("biz", biz), zap.String("email", email))

	// 生成验证码
	code := svc.generateCode()
	svc.logger.Debug("生成验证码成功", zap.String("code", code))

	// 塞进 redis
	err := svc.repo.Store(ctx, biz, email, code)
	if err != nil {
		svc.logger.Error("存储验证码失败", zap.Error(err), zap.String("biz", biz), zap.String("email", email))
		return err
	}
	svc.logger.Info("验证码存储成功", zap.String("biz", biz), zap.String("email", email))

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
