package email_domain

import (
	"strings"
	"time"
)

// EmailConfig 邮件服务器配置领域模型
// 对应数据库表EM01
// 存储SMTP邮件服务器配置信息
type EmailConfig struct {
	EMA001 string    // 配置ID（主键）
	EMA002 string    // 发送方邮箱地址
	EMA003 string    // 发送方邮箱授权码
	EMA004 string    // SMTP服务器地址
	EMA005 int       // SMTP服务器端口
	EMA006 string    // 配置名称（如：主邮箱、备用邮箱）
	EMA007 int       // 配置状态：1-启用，2-禁用
	EMA008 time.Time // 创建时间
	EMA009 time.Time // 修改时间
	EMA010 string    // 备注信息
}

// EmailTemplate 邮件模板领域模型
// 对应数据库表EM02
// 存储系统中使用的各类邮件模板
type EmailTemplate struct {
	EMB001 string    // 模板ID（主键）
	EMB002 string    // 模板名称
	EMB003 string    // 模板标识（唯一标识）
	EMB004 string    // 邮件主题
	EMB005 string    // 邮件内容（支持HTML格式，可包含占位符）
	EMB006 int       // 模板状态：1-启用，2-禁用
	EMB007 time.Time // 创建时间
	EMB008 time.Time // 修改时间
	EMB009 string    // 备注信息
}

// NewEmailConfig 创建新的邮件服务器配置
func NewEmailConfig(senderEmail, authPwd, smtpHost string, smtpPort int, configName string) *EmailConfig {
	return &EmailConfig{
		EMA002: senderEmail,
		EMA003: authPwd,
		EMA004: smtpHost,
		EMA005: smtpPort,
		EMA006: configName,
		EMA007: 1,          // 默认启用配置
		EMA008: time.Now(), // 创建时间
		EMA009: time.Now(), // 修改时间
	}
}

// NewEmailTemplate 创建新的邮件模板
func NewEmailTemplate(name, identifier, subject, content string) *EmailTemplate {
	return &EmailTemplate{
		EMB002: name,
		EMB003: identifier,
		EMB004: subject,
		EMB005: content,
		EMB006: 1, // 默认启用
	}
}

// IsEnabled 检查配置是否启用
func (ec *EmailConfig) IsEnabled() bool {
	return ec.EMA007 == 1
}

// Enable 启用配置
func (ec *EmailConfig) Enable() {
	ec.EMA007 = 1
	ec.EMA009 = time.Now() // 更新修改时间
}

// Disable 禁用配置
func (ec *EmailConfig) Disable() {
	ec.EMA007 = 2
	ec.EMA009 = time.Now() // 更新修改时间
}

// IsEnabled 检查模板是否启用
func (et *EmailTemplate) IsEnabled() bool {
	return et.EMB006 == 1
}

// Enable 启用模板
func (et *EmailTemplate) Enable() {
	et.EMB006 = 1
}

// Disable 禁用模板
func (et *EmailTemplate) Disable() {
	et.EMB006 = 2
}

// ReplacePlaceholders 替换邮件内容中的占位符
// placeholders: 占位符映射，键为占位符名称（不含{}），值为替换内容
func (et *EmailTemplate) ReplacePlaceholders(placeholders map[string]string) string {
	content := et.EMB005
	for key, value := range placeholders {
		// 替换形如 {key} 的占位符
		placeholder := "{" + key + "}"
		content = strings.ReplaceAll(content, placeholder, value)
	}
	return content
}
