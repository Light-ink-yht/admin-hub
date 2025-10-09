package user_domain

// SendSignupCodeRequest 发送注册验证码请求体
type SendSignupCodeRequest struct {
	AAA018 string // 前端输入框输入的值 (在后端判断是邮箱还是手机号注册)
}
