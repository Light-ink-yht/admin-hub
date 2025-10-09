package user_domain

type Login struct {
	AAA002 string // 邮箱 全局唯一 （前端不输入）
	AAA003 string // 手机号 全局唯一 （前端不输入）
	AAA004 string // 密码
	AAA017 string // 验证码 (captcha_code)
	AAA018 string // 前端输入框输入的值 (在后端判断是邮箱还是手机号注册)
	AAA019 string // captcha_id
}

// Validate 校验请求参数
func (req *Login) Validate() error {
	if req.AAA002 != "" {
		// 校验邮箱格式
		if match, _ := EmailRegex.MatchString(req.AAA002); !match {
			return ErrTheMailboxIsNotInTheRightFormat
		}
	}

	if req.AAA003 != "" {
		// 校验手机号格式
		if match, _ := PhoneRegex.MatchString(req.AAA003); !match {
			return ErrTheMobilePhoneNumberFormatIsInvalid
		}
	}

	// 校验密码格式
	if match, _ := PasswordRegex.MatchString(req.AAA004); !match {
		return ErrThePasswordIsNotInTheRightFormat
	}

	return nil
}
