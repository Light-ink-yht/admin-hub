package user_web

import (
	"fmt"
	"net/http"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/Light-ink-yht/admin-hub/internal/service/code_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/web"
	"github.com/Light-ink-yht/admin-hub/ioc/middleware"
	"github.com/Light-ink-yht/admin-hub/pkg/res"
	"github.com/gin-gonic/gin"
)

// biz 是一个常量，用于标识业务类型为 "signup"
const biz = "signup"

// 这是一个类型断言，用于在编译时检查 UserHandler 是否实现了 web.Handler 接口
var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
	codeSvc    code_svc.EmailServiceFace
	logService log_svc.LogService
}

func NewUserHandler(codeSvc code_svc.EmailServiceFace, logService log_svc.LogService) *UserHandler {
	return &UserHandler{
		codeSvc:    codeSvc,
		logService: logService,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/user")
	router.PUT("/signup/email/code", h.SendSignupEmailCode) // 发送邮箱验证码
	//router.POST("/signup/email/code", h.SignupEmail)        // 邮箱注册
	//router.POST("/login", h.Login)                          // 登录
	//router.POST("/layout", h.Layout)                        // 登出
	//router.GET("/me", h.GetMyInfo)                          // 获取当前用户信息
	//router.PUT("/me", h.EditMyInfo)                         // 编辑当前用户信息
	//router.PUT("/password", h.EditPassword)                 // 修改用户密码
}

// SendSignupEmailCode 发送邮箱验证码
func (h *UserHandler) SendSignupEmailCode(ctx *gin.Context) {
	var req user_domain.SendSignupEmailCodeRequest

	// 获取客户端IP（从中间件获取）
	clientIP, _ := middleware.GetRequestContextInfo(ctx)

	if err := ctx.Bind(&req); err != nil {

		// 记录业务错误日志
		content := fmt.Sprintf("请求绑定失败: %s", err.Error())
		h.logService.LogBusiness(ctx, log_domain.LogLevelError, "用户Web层", "发送注册邮箱验证码", "", "", "",
			content, "失败", err.Error(), clientIP, nil)
		return
	}

	body := `<p>您的验证码是：<strong>{code}</strong></p>
          <p>请在10分钟内使用该验证码完成验证。</p>`
	err := h.codeSvc.Send(ctx, biz, req.AAA002, body)

	if err != nil {
		// 获取客户端IP和处理时间

		// 记录业务错误日志
		content := fmt.Sprintf("发送验证码失败: %s", req.AAA002)
		h.logService.LogBusiness(ctx, log_domain.LogLevelError, "用户Web层", "发送注册邮箱验证码", "", "", "",
			content, "失败", err.Error(), clientIP, nil)
		ctx.JSON(http.StatusOK, res.Fail("系统异常"))
		return
	}

	// 记录业务成功日志
	content := fmt.Sprintf("发送验证码成功: %s", req.AAA002)
	h.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "用户Web层", "发送注册邮箱验证码", "", "", "",
		content, "成功", "", "", map[string]interface{}{
			"email": req.AAA002,
		})

	ctx.JSON(http.StatusOK, res.Success("验证码发送成功"))
}

func (h *UserHandler) SignupEmail(ctx *gin.Context) {

}
