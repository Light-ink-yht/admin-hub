package user_domain

// SendSignupEmailCodeRequest 发送注册验证码请求体
type SendSignupEmailCodeRequest struct {
	AAA002 string // 邮箱
}

// Validate 校验请求参数
func (req *SendSignupEmailCodeRequest) Validate() error {
	// 校验邮箱格式
	if match, _ := emailRegex.MatchString(req.AAA002); !match {
		return ErrTheMailboxIsNotInTheRightFormat
	}
	return nil
}
