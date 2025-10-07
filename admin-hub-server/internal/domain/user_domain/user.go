package user_domain

import (
	"errors"
	"time"

	"github.com/Light-ink-yht/admin-hub/pkg/model"
	"github.com/dlclark/regexp2"
)

// AA01 用户模型
type AA01 struct {
	model.LN01
	AAA001 int64     // 用户id
	AAA002 string    // 邮箱 全局唯一
	AAA003 string    // 手机号 全局唯一
	AAA004 string    // 密码
	AAA005 string    // 昵称
	AAA006 string    // 姓名
	AAA007 string    // 头像
	AAA008 string    // 性别 1 男，2 女，3 未知
	AAA009 time.Time // 生日
	AAA010 string    // 状态 1 启用，2 禁用
	AAA011 string    // 备注
	AAA012 int       // 登录次数
	AAA013 time.Time // 最后登录时间
	AAA014 string    // 最后登录ip
	AAA015 time.Time // 密码修改时间
}

var (
	emailRegex                             = regexp2.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, regexp2.None)
	passwordRegex                          = regexp2.MustCompile(`^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[._~!@#$^&*])[A-Za-z0-9._~!@#$^&*]{8,20}$`, regexp2.None)
	phoneRegex                             = regexp2.MustCompile(`^1[3-9]\d{9}$`, regexp2.None)
	ErrTheMailboxIsNotInTheRightFormat     = errors.New("电子邮件格式无效")
	ErrThePasswordIsNotInTheRightFormat    = errors.New("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符")
	ErrThePasswordIsInconsistentTwice      = errors.New("两次密码不一致")
	ErrTheMobilePhoneNumberFormatIsInvalid = errors.New("手机号格式无效")
)

// ValidateEmail 校验邮箱格式
// 入参：待校验的邮箱字符串
// 返回：nil（格式正确）或 ErrTheMailboxIsNotInTheRightFormat（格式错误）
func ValidateEmail(email string) error {
	// 先判断是否为空（空字符串直接判定为格式无效）
	if email == "" {
		return ErrTheMailboxIsNotInTheRightFormat
	}
	// regexp2 的 MatchString 方法返回 bool（是否匹配）和 error（正则本身是否有误，这里预编译过，error 可忽略）
	match, _ := emailRegex.MatchString(email)
	if !match {
		return ErrTheMailboxIsNotInTheRightFormat
	}
	return nil
}

// ValidatePhone 校验手机号格式
// 入参：待校验的手机号字符串
// 返回：nil（格式正确）或 ErrTheMobilePhoneNumberFormatIsInvalid（格式错误）
func ValidatePhone(phone string) error {
	if phone == "" {
		return ErrTheMobilePhoneNumberFormatIsInvalid
	}
	match, _ := phoneRegex.MatchString(phone)
	if !match {
		return ErrTheMobilePhoneNumberFormatIsInvalid
	}
	return nil
}

// JudgeInputType 判断输入是邮箱、手机号还是无效格式（类似前端的类型判断）
// 入参：待判断的字符串（邮箱/手机号）
// 返回："email"（邮箱）、"phone"（手机号）、""（无效格式）
func JudgeInputType(input string) (string, error) {
	// 先判断手机号（纯数字特征，避免与邮箱混淆）
	if _, err := phoneRegex.MatchString(input); err == nil {
		return "phone", err
	}
	// 再判断邮箱（含 @ 特征）
	if _, err := emailRegex.MatchString(input); err == nil {
		return "email", err
	}
	// 两者都不匹配，返回空字符串（无效格式）
	return "", nil
}
