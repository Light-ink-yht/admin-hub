package sms_domain

import (
	"strings"
	"time"
)

// SmsConfig 短信服务商配置领域模型
// 对应数据库表SM01
// 存储短信服务商的配置信息，如腾讯云、阿里云等不同服务商的密钥信息
type SmsConfig struct {
	SMA001 string    // 配置ID（主键）
	SMA002 string    // 配置名称（如：腾讯云主账号、阿里云备用账号）
	SMA003 string    // 短信服务商类型（tencent、aliyun等）
	SMA004 string    // 应用ID（AppId/SdkAppId）
	SMA005 string    // 应用密钥（SecretId/AccessKeyId）
	SMA006 string    // 应用密钥（SecretKey/AccessKeySecret）
	SMA007 string    // 短信签名
	SMA008 int       // 配置状态：1-启用，2-禁用
	SMA009 time.Time // 创建时间
	SMA010 time.Time // 修改时间
	SMA011 string    // 备注信息
}

// SmsTemplate 短信模板领域模型
// 对应数据库表SM02
// 存储短信模板信息，包括模板ID和模板内容描述
type SmsTemplate struct {
	SMB001 string    // 模板ID（主键）
	SMB002 string    // 模板名称
	SMB003 string    // 模板标识（唯一标识）
	SMB004 string    // 短信服务商模板ID
	SMB005 string    // 模板内容描述
	SMB006 int       // 模板状态：1-启用，2-禁用
	SMB007 time.Time // 创建时间
	SMB008 time.Time // 修改时间
	SMB009 string    // 备注信息
}

// NewSmsConfig 创建新的短信服务商配置
func NewSmsConfig(name, providerType, appId, secretId, secretKey, signName string) *SmsConfig {
	return &SmsConfig{
		SMA002: name,
		SMA003: providerType,
		SMA004: appId,
		SMA005: secretId,
		SMA006: secretKey,
		SMA007: signName,
		SMA008: 1,          // 默认启用配置
		SMA009: time.Now(), // 创建时间
		SMA010: time.Now(), // 修改时间
	}
}

// NewSmsTemplate 创建新的短信模板
func NewSmsTemplate(name, identifier, templateId, contentDesc string) *SmsTemplate {
	return &SmsTemplate{
		SMB002: name,
		SMB003: identifier,
		SMB004: templateId,
		SMB005: contentDesc,
		SMB006: 1, // 默认启用
	}
}

// IsEnabled 检查配置是否启用
func (sc *SmsConfig) IsEnabled() bool {
	return sc.SMA008 == 1
}

// Enable 启用配置
func (sc *SmsConfig) Enable() {
	sc.SMA008 = 1
	sc.SMA010 = time.Now() // 更新修改时间
}

// Disable 禁用配置
func (sc *SmsConfig) Disable() {
	sc.SMA008 = 2
	sc.SMA010 = time.Now() // 更新修改时间
}

// IsEnabled 检查模板是否启用
func (st *SmsTemplate) IsEnabled() bool {
	return st.SMB006 == 1
}

// Enable 启用模板
func (st *SmsTemplate) Enable() {
	st.SMB006 = 1
}

// Disable 禁用模板
func (st *SmsTemplate) Disable() {
	st.SMB006 = 2
}

// GetProviderType 获取短信服务商类型（小写）
func (sc *SmsConfig) GetProviderType() string {
	return strings.ToLower(sc.SMA003)
}
