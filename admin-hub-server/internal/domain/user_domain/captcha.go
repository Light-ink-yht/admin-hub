package user_domain

// CaptchaRequest 图形验证码请求结构体
type CaptchaRequest struct {
	// 验证码ID，用于验证时匹配
	CaptchaId string `json:"captcha_id"`
	// 用户输入的验证码
	CaptchaCode string `json:"captcha_code"`
}
