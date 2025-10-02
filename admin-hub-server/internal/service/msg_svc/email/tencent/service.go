package tencent

import (
	"context"
	"fmt"
	"github.com/go-gomail/gomail"
)

type EmailService struct {
	authEmail string // 发送方的邮箱地址
	authPwd   string // 发送方的邮箱授权码（非QQ密码）
	smtpHost  string // SMTP服务器地址
	smtpPort  int    // SMTP服务器端口
}

func NewEmailService(authEmail, authPwd, smtpHost string, smtpPort int) *EmailService {
	return &EmailService{
		authEmail: authEmail,
		authPwd:   authPwd,
		smtpHost:  smtpHost,
		smtpPort:  smtpPort,
	}
}

func (e *EmailService) Send(ctx context.Context, email string, subject, body string) error {
	// 创建一个新的邮件发送器
	client := gomail.NewDialer(e.smtpHost, e.smtpPort, e.authEmail, e.authPwd)
	// 创建一个新的邮件消息
	m := gomail.NewMessage()
	// 设置邮件的发件人
	m.SetHeader("From", e.authEmail)
	// 设置邮件的收件人
	m.SetHeader("To", email)
	// 设置邮件的主题
	m.SetHeader("Subject", subject)
	// 设置邮件的正文
	m.SetBody("text/html", body)
	// 发送邮件
	err := client.DialAndSend(m)
	if err != nil {
		return fmt.Errorf(" %s 发送邮件失败", email)
	}
	return nil
}
