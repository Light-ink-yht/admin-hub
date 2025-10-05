package email_web

import (
	"net/http"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/email_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/web"
	"github.com/Light-ink-yht/admin-hub/ioc/middleware"
	"github.com/Light-ink-yht/admin-hub/pkg/res"
	"github.com/gin-gonic/gin"
)

// 确保EmailHandler实现了web.Handler接口
var _ web.Handler = (*EmailHandler)(nil)

type EmailHandler struct {
	emailConfigRepo   email_repo.EmailConfigRepository
	emailTemplateRepo email_repo.EmailTemplateRepository
	logService        log_svc.LogService
}

func NewEmailHandler(
	emailConfigRepo email_repo.EmailConfigRepository,
	emailTemplateRepo email_repo.EmailTemplateRepository,
	logService log_svc.LogService,
) *EmailHandler {
	return &EmailHandler{
		emailConfigRepo:   emailConfigRepo,
		emailTemplateRepo: emailTemplateRepo,
		logService:        logService,
	}
}

func (h *EmailHandler) RegisterRoutes(r *gin.RouterGroup) {
	emailRouter := r.Group("/email")
	{
		// 邮件服务器配置相关接口
		configRouter := emailRouter.Group("/config")
		{
			configRouter.POST("", h.CreateEmailConfig)                 // 创建邮件配置
			configRouter.PUT("/:id", h.UpdateEmailConfig)              // 更新邮件配置
			configRouter.DELETE("/:id", h.DeleteEmailConfig)           // 删除邮件配置
			configRouter.GET("", h.GetEmailConfigs)                    // 获取邮件配置（支持ID查询和列表查询）
			configRouter.PUT("/:id/status", h.UpdateEmailConfigStatus) // 更新邮件配置状态
			configRouter.GET("/default", h.GetDefaultEmailConfig)      // 获取默认邮件配置
			configRouter.PUT("/:id/default", h.SetDefaultEmailConfig)  // 设置默认邮件配置
		}

		// 邮件模板相关接口
		templateRouter := emailRouter.Group("/template")
		{
			templateRouter.POST("", h.CreateEmailTemplate)                 // 创建邮件模板
			templateRouter.PUT("/:id", h.UpdateEmailTemplate)              // 更新邮件模板
			templateRouter.DELETE("/:id", h.DeleteEmailTemplate)           // 删除邮件模板
			templateRouter.GET("", h.GetEmailTemplates)                    // 获取邮件模板（支持ID、类型查询和列表查询）
			templateRouter.PUT("/:id/status", h.UpdateEmailTemplateStatus) // 更新邮件模板状态
			templateRouter.GET("/search", h.SearchEmailTemplates)          // 搜索邮件模板
		}
	}
}

// CreateEmailConfig 创建邮件配置
func (h *EmailHandler) CreateEmailConfig(ctx *gin.Context) {
	var req email_domain.EmailConfig
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "创建邮件配置", "请求参数绑定失败", err)
		return
	}

	// 设置默认值
	if req.EMA007 == 0 {
		req.EMA007 = 1 // 默认启用
	}

	if err := h.emailConfigRepo.Create(&req); err != nil {
		h.logError(ctx, "创建邮件配置", "保存失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "创建邮件配置", "成功", map[string]interface{}{
		"config_id": req.EMA001,
	})
	ctx.JSON(http.StatusOK, res.Success("创建成功"))
}

// UpdateEmailConfig 更新邮件配置
func (h *EmailHandler) UpdateEmailConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	var req email_domain.EmailConfig
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "更新邮件配置", "请求参数绑定失败", err)
		return
	}

	req.EMA001 = id
	if err := h.emailConfigRepo.Update(&req); err != nil {
		h.logError(ctx, "更新邮件配置", "保存失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新邮件配置", "成功", map[string]interface{}{
		"config_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("更新成功"))
}

// DeleteEmailConfig 删除邮件配置
func (h *EmailHandler) DeleteEmailConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.emailConfigRepo.Delete(id); err != nil {
		h.logError(ctx, "删除邮件配置", "失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "删除邮件配置", "成功", map[string]interface{}{
		"config_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("删除成功"))
}

// GetEmailConfigs 获取邮件配置
// 支持通过ID查询单个配置或获取所有配置
func (h *EmailHandler) GetEmailConfigs(ctx *gin.Context) {
	id := ctx.Query("id")

	// 如果提供了ID，获取单个配置
	if id != "" {
		config, err := h.emailConfigRepo.GetByID(id)
		if err != nil {
			h.logError(ctx, "获取邮件配置", "查询失败", err)
			ctx.JSON(http.StatusOK, res.Fail(err.Error()))
			return
		}

		if config == nil {
			h.logError(ctx, "获取邮件配置", "配置不存在", nil)
			ctx.JSON(http.StatusOK, res.Fail("配置不存在"))
			return
		}

		h.logSuccess(ctx, "获取邮件配置", "成功", map[string]interface{}{
			"config_id": id,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取邮件配置成功", config))
		return
	}

	// 否则获取所有配置
	configs, err := h.emailConfigRepo.GetAll()
	if err != nil {
		h.logError(ctx, "获取所有邮件配置", "查询失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "获取所有邮件配置", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取所有邮件配置成功", configs))
}

// UpdateEmailConfigStatus 更新邮件配置状态
func (h *EmailHandler) UpdateEmailConfigStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	// 验证状态参数
	if status != "enable" && status != "disable" {
		h.logError(ctx, "更新邮件配置状态", "状态参数无效", nil)
		ctx.JSON(http.StatusOK, res.Fail("状态参数无效，应为'enable'或'disable'"))
		return
	}

	var err error
	if status == "enable" {
		err = h.emailConfigRepo.Enable(id)
	} else {
		err = h.emailConfigRepo.Disable(id)
	}

	if err != nil {
		h.logError(ctx, "更新邮件配置状态", "失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新邮件配置状态", "成功", map[string]interface{}{
		"config_id": id,
		"status":    status,
	})
	ctx.JSON(http.StatusOK, res.Success("状态更新成功"))
}

// GetDefaultEmailConfig 获取默认邮件配置
func (h *EmailHandler) GetDefaultEmailConfig(ctx *gin.Context) {
	config, err := h.emailConfigRepo.GetDefault()
	if err != nil {
		h.logError(ctx, "获取默认邮件配置", "查询失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	if config == nil {
		h.logError(ctx, "获取默认邮件配置", "默认配置不存在", nil)
		ctx.JSON(http.StatusOK, res.Fail("默认配置不存在"))
		return
	}

	h.logSuccess(ctx, "获取默认邮件配置", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取默认邮件配置成功", config))
}

// SetDefaultEmailConfig 设置默认邮件配置
func (h *EmailHandler) SetDefaultEmailConfig(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.emailConfigRepo.SetDefault(id); err != nil {
		h.logError(ctx, "设置默认邮件配置", "设置失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "设置默认邮件配置", "成功", map[string]interface{}{
		"config_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("设置默认配置成功"))
}

// GetEmailTemplates 获取邮件模板
// 支持通过ID查询单个模板、通过类型查询模板或获取所有模板
func (h *EmailHandler) GetEmailTemplates(ctx *gin.Context) {
	id := ctx.Query("id")
	templateType := ctx.Query("type")

	// 如果提供了ID，获取单个模板
	if id != "" {
		template, err := h.emailTemplateRepo.GetByID(id)
		if err != nil {
			h.logError(ctx, "获取邮件模板", "查询失败", err)
			ctx.JSON(http.StatusOK, res.Fail(err.Error()))
			return
		}

		if template == nil {
			h.logError(ctx, "获取邮件模板", "模板不存在", nil)
			ctx.JSON(http.StatusOK, res.Fail("模板不存在"))
			return
		}

		h.logSuccess(ctx, "获取邮件模板", "成功", map[string]interface{}{
			"template_id": id,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取邮件模板成功", template))
		return
	}

	// 如果提供了类型，根据类型获取模板
	if templateType != "" {
		template, err := h.emailTemplateRepo.GetByType(templateType)
		if err != nil {
			h.logError(ctx, "获取邮件模板", "查询失败", err)
			ctx.JSON(http.StatusOK, res.Fail(err.Error()))
			return
		}

		if template == nil {
			h.logError(ctx, "获取邮件模板", "模板不存在", nil)
			ctx.JSON(http.StatusOK, res.Fail("模板不存在"))
			return
		}

		h.logSuccess(ctx, "获取邮件模板", "成功", map[string]interface{}{
			"template_type": templateType,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取邮件模板成功", template))
		return
	}

	// 否则获取所有模板
	templates, err := h.emailTemplateRepo.GetAll()
	if err != nil {
		h.logError(ctx, "获取所有邮件模板", "查询失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "获取所有邮件模板", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取所有邮件模板成功", templates))
}

// UpdateEmailTemplateStatus 更新邮件模板状态
func (h *EmailHandler) UpdateEmailTemplateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	// 验证状态参数
	if status != "enable" && status != "disable" {
		h.logError(ctx, "更新邮件模板状态", "状态参数无效", nil)
		ctx.JSON(http.StatusOK, res.Fail("状态参数无效，应为'enable'或'disable'"))
		return
	}

	var err error
	if status == "enable" {
		err = h.emailTemplateRepo.Enable(id)
	} else {
		err = h.emailTemplateRepo.Disable(id)
	}

	if err != nil {
		h.logError(ctx, "更新邮件模板状态", "失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新邮件模板状态", "成功", map[string]interface{}{
		"template_id": id,
		"status":      status,
	})
	ctx.JSON(http.StatusOK, res.Success("状态更新成功"))
}

// CreateEmailTemplate 创建邮件模板
func (h *EmailHandler) CreateEmailTemplate(ctx *gin.Context) {
	var req email_domain.EmailTemplate
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "创建邮件模板", "请求参数绑定失败", err)
		return
	}

	// 设置默认值
	if req.EMB008.IsZero() {
		req.EMB008 = time.Now() // 设置当前时间为更新时间
	}

	// 设置创建时间
	req.EMB007 = time.Now()

	if err := h.emailTemplateRepo.Create(&req); err != nil {
		h.logError(ctx, "创建邮件模板", "保存失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "创建邮件模板", "成功", map[string]interface{}{
		"template_id": req.EMB001,
	})
	ctx.JSON(http.StatusOK, res.Success("创建成功"))
}

// UpdateEmailTemplate 更新邮件模板
func (h *EmailHandler) UpdateEmailTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	var req email_domain.EmailTemplate
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "更新邮件模板", "请求参数绑定失败", err)
		return
	}

	// 设置更新时间
	req.EMB008 = time.Now()

	req.EMB001 = id
	if err := h.emailTemplateRepo.Update(&req); err != nil {
		h.logError(ctx, "更新邮件模板", "保存失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新邮件模板", "成功", map[string]interface{}{
		"template_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("更新成功"))
}

// DeleteEmailTemplate 删除邮件模板
func (h *EmailHandler) DeleteEmailTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.emailTemplateRepo.Delete(id); err != nil {
		h.logError(ctx, "删除邮件模板", "失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "删除邮件模板", "成功", map[string]interface{}{
		"template_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("删除成功"))
}

// SearchEmailTemplates 搜索邮件模板
func (h *EmailHandler) SearchEmailTemplates(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	templates, err := h.emailTemplateRepo.Search(keyword)
	if err != nil {
		h.logError(ctx, "搜索邮件模板", "查询失败", err)
		ctx.JSON(http.StatusOK, res.Fail(err.Error()))
		return
	}

	h.logSuccess(ctx, "搜索邮件模板", "成功", map[string]interface{}{
		"keyword": keyword,
	})
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取成功", templates))
}

// 记录错误日志
func (h *EmailHandler) logError(ctx *gin.Context, module, content string, err error) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	logError := ""
	if err != nil {
		logError = err.Error()
	}
	h.logService.LogBusiness(ctx, log_domain.LogLevelError, "邮件Web层", module, "", "", "",
		content, "失败", logError, clientIP, nil)
}

// 记录成功日志
func (h *EmailHandler) logSuccess(ctx *gin.Context, module, content string, extra map[string]interface{}) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	h.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "邮件Web层", module, "", "", "",
		content, "成功", "", clientIP, extra)
}
