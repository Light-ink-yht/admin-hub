package code_domain

import (
	"time"
)

// Code 验证码领域模型
type Code struct {
	CDA001 string    // 主键ID（业务主键）
	CDA002 string    // 业务类型(biz)，如email、sms等
	CDA003 string    // 用户输入(input)，如邮箱地址、手机号等
	CDA004 string    // 验证码内容
	CDA005 time.Time // 创建时间
	CDA006 time.Time // 过期时间
	CDA007 int       // 验证次数限制
	CDA008 int       // 已验证次数
	CDA009 int       // 状态(1:有效, 2:已使用, 3:已过期, 4:已失效)
	CDA010 string    // 备注信息
}

// NewCode 创建一个新的验证码模型
func NewCode(biz, input, code string, expireAt time.Time, maxAttempts int) *Code {
	return &Code{
		CDA002: biz,
		CDA003: input,
		CDA004: code,
		CDA005: time.Now(),
		CDA006: expireAt,
		CDA007: maxAttempts,
		CDA008: 0,
		CDA009: 1, // 默认为有效状态
	}
}

// IsExpired 检查验证码是否已过期
func (c *Code) IsExpired() bool {
	return time.Now().After(c.CDA006)
}

// CanVerify 检查验证码是否可以继续验证
func (c *Code) CanVerify() bool {
	return c.CDA009 == 1 && !c.IsExpired() && c.CDA008 < c.CDA007
}

// VerifySuccess 验证成功后更新状态
func (c *Code) VerifySuccess() {
	c.CDA009 = 2 // 已使用
	c.CDA008++
}

// VerifyFailed 验证失败后更新状态
func (c *Code) VerifyFailed() {
	c.CDA008++
	if c.CDA008 >= c.CDA007 {
		c.CDA009 = 4 // 已失效（验证次数过多）
	}
}

// Expire 使验证码过期
func (c *Code) Expire() {
	c.CDA009 = 3 // 已过期
}
