package email

import "context"

// EmailService 定义了一个邮件服务接口
type EmailService interface {
	// Send 发送邮件
	Send(ctx context.Context, email string, subject, body string) error
}
