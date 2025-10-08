package user_web

import (
	"errors"
	"net/http"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/Light-ink-yht/admin-hub/internal/service/code_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/user_svc"
	"github.com/Light-ink-yht/admin-hub/internal/web"
	"github.com/Light-ink-yht/admin-hub/internal/web/middleware"
	"github.com/Light-ink-yht/admin-hub/pkg/res"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// biz 是一个常量，用于标识业务类型为 "signup"
const biz = "signup"

// 这是一个类型断言，用于在编译时检查 UserHandler 是否实现了 web.Handler 接口
var _ web.Handler = (*UserHandler)(nil)

// UserHandler 用户处理器
type UserHandler struct {
	svc          user_svc.UserService
	emailSvc     code_svc.EmailService
	smsSvc       code_svc.SmsServiceFace
	logService   log_svc.LogService
	captchaStore base64Captcha.Store
}

func NewUserHandler(svc user_svc.UserService, emailSvc code_svc.EmailService, smsSvc code_svc.SmsServiceFace, logService log_svc.LogService) *UserHandler {
	return &UserHandler{
		svc:          svc,
		emailSvc:     emailSvc,
		smsSvc:       smsSvc,
		logService:   logService,
		captchaStore: base64Captcha.DefaultMemStore,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/user")
	router.PUT("/signup/code", h.SendSignupCode) // 发送注册验证码
	router.POST("/signup/code", h.Signup)        // 邮箱/手机号注册
	router.GET("/captcha", h.GenerateCaptcha)    // 生成图形验证码
	router.POST("/login", h.Login)               // 登录
	//router.POST("/layout", h.Layout)           // 登出
	//router.GET("/me", h.GetMyInfo)             // 获取当前用户信息
	//router.PUT("/me", h.EditMyInfo)            // 编辑当前用户信息
	//router.PUT("/password", h.EditPassword)    // 修改用户密码
}

// SendSignupCode 发送注册验证码
func (h *UserHandler) SendSignupCode(ctx *gin.Context) {
	var req user_domain.SendSignupCodeRequest

	if err := ctx.Bind(&req); err != nil {
		h.logError(ctx, "发送注册邮箱验证码", "请求绑定失败", err)
		return
	}

	// 判断前端传过来的 AAA018 是邮箱还是手机号
	inputType := user_domain.JudgeInputType(req.AAA018)

	if inputType == "email" {
		err := h.emailSvc.Send(ctx, biz, req.AAA018, "发送邮箱验证码")

		if errors.Is(err, user_domain.ErrTheMailboxIsNotInTheRightFormat) {
			// 如果邮箱格式无效，返回错误信息
			h.logWarn(ctx, "用户邮箱注册验证码", "电子邮件格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("电子邮件格式无效"))
			return
		}

		if err != nil {
			h.logError(ctx, "发送注册邮箱验证码", "发送验证码失败", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		h.logSuccess(ctx, "发送注册邮箱验证码", "发送验证码成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("验证码发送成功"))
	} else if inputType == "phone" {
		err := h.smsSvc.Send(ctx, biz, req.AAA018)
		if errors.Is(err, user_domain.ErrTheMobilePhoneNumberFormatIsInvalid) {
			// 如果手机号格式无效，返回错误信息
			h.logWarn(ctx, "用户手机号注册", "手机号格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("手机号格式无效"))
			return
		}

		if err != nil {
			h.logError(ctx, "发送注册手机验证码", "发送验证码失败", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		h.logSuccess(ctx, "发送注册手机验证码", "发送验证码成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("验证码发送成功"))
	}

}

// Signup 注册
func (h *UserHandler) Signup(ctx *gin.Context) {
	var req user_domain.SignUp

	if err := ctx.Bind(&req); err != nil {
		h.logError(ctx, "注册", "请求绑定失败", err)
		return
	}

	// 判断前端传过来的 AAA018 是邮箱还是手机号
	inputType := user_domain.JudgeInputType(req.AAA018)

	if inputType == "email" {
		req.AAA002 = req.AAA018
		// 验证验证码
		ok, err := h.emailSvc.Verify(ctx, biz, req.AAA002, req.AAA017)
		if errors.Is(err, code_svc.ErrCodeSendTooMany) {
			// 如果发送验证码太频繁，返回错误信息
			h.logWarn(ctx, "验证验证码", "发送验证码太频繁", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("发送验证码太频繁"))
			return
		}
		if errors.Is(err, code_svc.RrrCodeVerifyTooMany) {
			// 如果验证次数太多，返回错误信息
			h.logWarn(ctx, "验证验证码", "验证次数太多", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("验证次数太多"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "验证验证码", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		if !ok {
			// 如果验证码错误，返回错误信息
			h.logWarn(ctx, "验证验证码", "验证码错误", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("验证码错误"))
			return
		}
		err = h.svc.Signup(ctx, &req)
		if errors.Is(err, user_domain.ErrTheMailboxIsNotInTheRightFormat) {
			// 如果邮箱格式无效，返回错误信息
			h.logWarn(ctx, "用户邮箱注册", "电子邮件格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("电子邮件格式无效"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
			// 如果密码格式不对，返回错误信息
			h.logWarn(ctx, "用户邮箱注册", "密码格式不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsInconsistentTwice) {
			// 如果两次密码不一致，返回错误信息
			h.logWarn(ctx, "用户邮箱注册", "两次密码不一致", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("两次密码不一致"))
			return
		}
		if errors.Is(err, user_svc.ErrUserDuplicateEmailOrPhone) {
			// 如果邮箱冲突，返回错误信息
			h.logWarn(ctx, "用户邮箱注册", "邮箱已存在", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("邮箱已存在"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "用户邮箱注册", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		h.logSuccess(ctx, "邮箱注册", "邮箱注册成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("注册成功"))
	} else if inputType == "phone" {
		req.AAA003 = req.AAA018
		// 验证验证码
		ok, err := h.smsSvc.Verify(ctx, biz, req.AAA003, req.AAA017)
		if errors.Is(err, code_svc.ErrCodeSendTooMany) {
			// 如果发送验证码太频繁，返回错误信息
			h.logWarn(ctx, "验证验证码", "发送验证码太频繁", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("发送验证码太频繁"))
			return
		}
		if errors.Is(err, code_svc.RrrCodeVerifyTooMany) {
			// 如果验证次数太多，返回错误信息
			h.logWarn(ctx, "验证验证码", "验证次数太多", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("验证次数太多"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "验证验证码", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		if !ok {
			// 如果验证码错误，返回错误信息
			h.logWarn(ctx, "验证验证码", "验证码错误", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("验证码错误"))
			return
		}
		err = h.svc.Signup(ctx, &req)
		if errors.Is(err, user_domain.ErrTheMobilePhoneNumberFormatIsInvalid) {
			// 如果邮箱格式无效，返回错误信息
			h.logWarn(ctx, "用户手机号注册", "手机号格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("手机号格式无效"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
			// 如果密码格式不对，返回错误信息
			h.logWarn(ctx, "用户手机号注册", "密码格式不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsInconsistentTwice) {
			// 如果两次密码不一致，返回错误信息
			h.logWarn(ctx, "用户手机号注册", "两次密码不一致", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("两次密码不一致"))
			return
		}
		if errors.Is(err, user_svc.ErrUserDuplicateEmailOrPhone) {
			// 如果邮箱冲突，返回错误信息
			h.logWarn(ctx, "用户手机号注册", "手机号已存在", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("手机号已存在"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "用户手机号注册", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		h.logSuccess(ctx, "用户手机号注册", "手机号注册成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("注册成功"))
	}
}

// Login 登录
func (h *UserHandler) Login(ctx *gin.Context) {
	var req user_domain.Login
	if err := ctx.Bind(&req); err != nil {
		h.logError(ctx, "登录", "请求绑定失败", err)
		return
	}

	// 验证图形验证码
	if !h.VerifyCaptcha(ctx, req.AAA019, req.AAA017) {
		h.logWarn(ctx, "登录", "图形验证码错误", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("图形验证码错误"))
		return
	}

	// 判断前端传过来的 AAA018 是邮箱还是手机号
	inputType := user_domain.JudgeInputType(req.AAA018)

	// 从Gin上下文中获取UserAgent
	userAgent := ctx.Request.UserAgent()

	if inputType == "email" {
		req.AAA002 = req.AAA018

		token, err := h.svc.EmailLogin(ctx, &req, userAgent)
		if errors.Is(err, user_domain.ErrTheMailboxIsNotInTheRightFormat) {
			// 如果邮箱格式无效，返回错误信息
			h.logWarn(ctx, "用户邮箱登录", "电子邮件格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("电子邮件格式无效"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
			// 如果密码格式不对，返回错误信息
			h.logWarn(ctx, "用户邮箱登录", "密码格式不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符"))
			return
		}
		if errors.Is(err, user_svc.ErrInvalidEmailOrPassword) {
			// 如果邮箱或密码错误，返回错误信息
			h.logWarn(ctx, "用户邮箱登录", "邮箱或密码不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("邮箱或密码不对"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "用户邮箱登录", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		// 设置响应头中的x-jwt-token
		ctx.Header("x-jwt-token", token)

		h.logSuccess(ctx, "用户邮箱登录", "邮箱登录成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("登录成功"))
	} else if inputType == "phone" {
		req.AAA003 = req.AAA018

		token, err := h.svc.PhoneLogin(ctx, &req, userAgent)
		if errors.Is(err, user_domain.ErrTheMobilePhoneNumberFormatIsInvalid) {
			// 如果邮箱格式无效，返回错误信息
			h.logWarn(ctx, "用户手机号登录", "手机号格式无效", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("手机号格式无效"))
			return
		}
		if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
			// 如果密码格式不对，返回错误信息
			h.logWarn(ctx, "用户手机号登录", "密码格式不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符"))
			return
		}
		if errors.Is(err, user_svc.ErrInvalidPhoneOrPassword) {
			// 如果邮箱或密码错误，返回错误信息
			h.logWarn(ctx, "用户手机号登录", "手机号或密码不对", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn("手机号或密码不对"))
			return
		}
		if err != nil {
			// 如果系统错误，返回错误信息
			h.logError(ctx, "用户手机号登录", "系统异常", err)
			ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
			return
		}

		// 设置响应头中的x-jwt-token
		ctx.Header("x-jwt-token", token)

		h.logSuccess(ctx, "用户手机号登录", "手机号登录成功", map[string]interface{}{})
		ctx.JSON(http.StatusOK, res.Success("登录成功"))
	}

}

// GenerateCaptcha 生成图形验证码
func (h *UserHandler) GenerateCaptcha(ctx *gin.Context) {
	// 配置验证码
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	// 创建验证码
	c := base64Captcha.NewCaptcha(driver, h.captchaStore)
	// 生成验证码
	id, b64s, _, err := c.Generate()
	if err != nil {
		h.logError(ctx, "生成图形验证码", "生成验证码失败", err)
		ctx.JSON(http.StatusOK, res.FailWithError("系统异常"))
		return
	}

	h.logSuccess(ctx, "生成图形验证码", "生成验证码失败", map[string]interface{}{})

	// 返回验证码信息
	ctx.JSON(http.StatusOK, res.SuccessWithData("图形验证码发送成功", map[string]string{
		"captcha_id":   id,
		"captcha_code": b64s,
	}))
}

// VerifyCaptcha 验证图形验证码
func (h *UserHandler) VerifyCaptcha(ctx *gin.Context, captchaId, captchaCode string) bool {
	// 验证验证码
	return h.captchaStore.Verify(captchaId, captchaCode, true)
}

// 记录错误日志
func (h *UserHandler) logError(ctx *gin.Context, module, content string, err error) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	logError := ""
	if err != nil {
		logError = err.Error()
	}
	h.logService.LogBusiness(ctx, log_domain.LogLevelError, "用户Web层", module, "", "", "",
		content, "失败", logError, clientIP, nil)
}

// 记录警告日志
func (h *UserHandler) logWarn(ctx *gin.Context, module, content string, err error) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	logError := ""
	if err != nil {
		logError = err.Error()
	}
	h.logService.LogBusiness(ctx, log_domain.LogLevelWarn, "用户Web层", module, "", "", "",
		content, "警告", logError, clientIP, nil)
}

// 记录成功日志
func (h *UserHandler) logSuccess(ctx *gin.Context, module, content string, extra map[string]interface{}) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	h.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "用户Web层", module, "", "", "",
		content, "成功", "", clientIP, extra)
}
