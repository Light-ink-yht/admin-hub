package user_domain

import "time"

type SignUp struct {
	AAA002 string // 邮箱 全局唯一 （前端不输入）
	AAA003 string // 手机号 全局唯一 （前端不输入）
	AAA004 string // 密码
	AAA016 string // 确认密码
	AAA017 string // 验证码
	AAA018 string // 前端输入框输入的值 (在后端判断是邮箱还是手机号注册)
}

// Validate 校验请求参数
func (req *SignUp) Validate() error {
	if req.AAA002 != "" {
		// 校验邮箱格式
		if match, _ := emailRegex.MatchString(req.AAA002); !match {
			return ErrTheMailboxIsNotInTheRightFormat
		}
	}

	if req.AAA003 != "" {
		// 校验手机号格式
		if match, _ := phoneRegex.MatchString(req.AAA003); !match {
			return ErrTheMobilePhoneNumberFormatIsInvalid
		}
	}

	// 校验密码格式
	if match, _ := passwordRegex.MatchString(req.AAA004); !match {
		return ErrThePasswordIsNotInTheRightFormat
	}

	// 确认密码是否一致
	if req.AAA004 != req.AAA016 {
		return ErrThePasswordIsInconsistentTwice
	}

	return nil
}

// DefaultUser 初始化默认用户
func DefaultUser(email, password, phone string, userId int64) *AA01 {
	return &AA01{
		AAA001: userId,
		AAA002: email,
		AAA003: phone,
		AAA004: password,
		AAA005: "昵称",
		AAA006: "",
		AAA007: "/uploads/avatar/logo.png",
		AAA008: "1",
		AAA009: time.Date(2025, time.October, 6, 0, 0, 0, 0, time.UTC),
		AAA010: "1",
		AAA011: "",
		AAA012: 0,
		AAA013: time.Time{},
		AAA014: "",
		AAA015: time.Time{},
	}
}
