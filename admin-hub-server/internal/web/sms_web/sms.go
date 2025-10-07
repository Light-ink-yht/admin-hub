package sms_web

import (
	"net/http"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/domain/sms_domain"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/sms_svc"
	"github.com/Light-ink-yht/admin-hub/internal/web"
	"github.com/Light-ink-yht/admin-hub/ioc/middleware"
	"github.com/Light-ink-yht/admin-hub/pkg/res"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// 确保SmsHandler实现了web.Handler接口
var _ web.Handler = (*SmsHandler)(nil)

type SmsHandler struct {
	smsService *sms_svc.Service
	logService log_svc.LogService
}

func NewSmsHandler(
	smsService *sms_svc.Service,
	logService log_svc.LogService,
) *SmsHandler {
	return &SmsHandler{
		smsService: smsService,
		logService: logService,
	}
}

func (h *SmsHandler) RegisterRoutes(r *gin.RouterGroup) {
	smsRouter := r.Group("/sms")
	{
		// 短信服务商配置相关接口
		configRouter := smsRouter.Group("/config")
		{
			configRouter.POST("", h.CreateSmsConfig)                 // 创建短信配置
			configRouter.PUT("/:id", h.UpdateSmsConfig)              // 更新短信配置
			configRouter.DELETE("/:id", h.DeleteSmsConfig)           // 删除短信配置
			configRouter.GET("", h.GetSmsConfigs)                    // 获取短信配置（支持ID查询和列表查询）
			configRouter.PUT("/:id/status", h.UpdateSmsConfigStatus) // 更新短信配置状态
			configRouter.GET("/default", h.GetDefaultSmsConfig)      // 获取默认短信配置
		}

		// 短信模板相关接口
		templateRouter := smsRouter.Group("/template")
		{
			templateRouter.POST("", h.CreateSmsTemplate)                 // 创建短信模板
			templateRouter.PUT("/:id", h.UpdateSmsTemplate)              // 更新短信模板
			templateRouter.DELETE("/:id", h.DeleteSmsTemplate)           // 删除短信模板
			templateRouter.GET("", h.GetSmsTemplates)                    // 获取短信模板（支持ID、标识查询和列表查询）
			templateRouter.PUT("/:id/status", h.UpdateSmsTemplateStatus) // 更新短信模板状态
			templateRouter.GET("/search", h.SearchSmsTemplates)          // 搜索短信模板
		}
	}
}

// CreateSmsConfig 创建短信配置
func (h *SmsHandler) CreateSmsConfig(ctx *gin.Context) {
	var req sms_domain.SmsConfig
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "创建短信配置", "请求参数绑定失败", err)
		return
	}

	// 设置默认值
	if req.SMA008 == 0 {
		req.SMA008 = 1 // 默认启用
	}

	// 设置创建和修改时间
	if req.SMA009.IsZero() {
		req.SMA009 = time.Now() // 创建时间
	}
	req.SMA010 = time.Now() // 修改时间

	if err := h.smsService.CreateSmsConfig(context.Background(), &req); err != nil {
		h.logError(ctx, "创建短信配置", "保存失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "创建短信配置", "成功", map[string]interface{}{
		"config_id": req.SMA001,
	})
	ctx.JSON(http.StatusOK, res.Success("创建成功"))
}

// UpdateSmsConfig 更新短信配置
func (h *SmsHandler) UpdateSmsConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	var req sms_domain.SmsConfig
	if err := ctx.BindJSON(&req); err != nil {
		h.logError(ctx, "更新短信配置", "请求参数绑定失败", err)
		return
	}

	// 设置ID和修改时间
	req.SMA001 = id
	req.SMA010 = time.Now() // 修改时间

	if err := h.smsService.UpdateSmsConfig(context.Background(), &req); err != nil {
		h.logError(ctx, "更新短信配置", "保存失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新短信配置", "成功", map[string]interface{}{
		"config_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("更新成功"))
}

// DeleteSmsConfig 删除短信配置
func (h *SmsHandler) DeleteSmsConfig(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := h.smsService.DeleteSmsConfig(context.Background(), id); err != nil {
		h.logError(ctx, "删除短信配置", "失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "删除短信配置", "成功", map[string]interface{}{
		"config_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("删除成功"))
}

// GetSmsConfigs 获取短信配置
// 支持通过ID查询单个配置或获取所有配置
func (h *SmsHandler) GetSmsConfigs(ctx *gin.Context) {
	id := ctx.Query("id")

	// 如果提供了ID，获取单个配置
	if id != "" {
		config, err := h.smsService.GetSmsConfig(context.Background(), id)
		if err != nil {
			h.logError(ctx, "获取短信配置", "查询失败", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
			return
		}

		if config == nil {
			h.logError(ctx, "获取短信配置", "配置不存在", nil)
			ctx.JSON(http.StatusOK, res.FailWithWarn("配置不存在"))
			return
		}

		h.logSuccess(ctx, "获取短信配置", "成功", map[string]interface{}{
			"config_id": id,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取短信配置成功", config))
		return
	}

	// 否则获取所有配置
	configs, err := h.smsService.GetAllSmsConfigs(context.Background())
	if err != nil {
		h.logError(ctx, "获取所有短信配置", "查询失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "获取所有短信配置", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取所有短信配置成功", configs))
}

// UpdateSmsConfigStatus 更新短信配置状态
func (h *SmsHandler) UpdateSmsConfigStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	// 验证状态参数
	if status != "enable" && status != "disable" {
		h.logError(ctx, "更新短信配置状态", "状态参数无效", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("状态参数无效，应为'enable'或'disable'"))
		return
	}

	// 获取当前配置
	config, err := h.smsService.GetSmsConfig(context.Background(), id)
	if err != nil {
		h.logError(ctx, "更新短信配置状态", "获取配置失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	if config == nil {
		h.logError(ctx, "更新短信配置状态", "配置不存在", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("配置不存在"))
		return
	}

	// 根据状态更新配置
	if status == "enable" {
		config.Enable()
	} else {
		config.Disable()
	}

	// 更新配置
	if err := h.smsService.UpdateSmsConfig(context.Background(), config); err != nil {
		h.logError(ctx, "更新短信配置状态", "更新配置失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新短信配置状态", "成功", map[string]interface{}{
		"config_id": id,
		"status":    status,
	})
	ctx.JSON(http.StatusOK, res.Success("状态更新成功"))
}

// GetDefaultSmsConfig 获取默认短信配置
func (h *SmsHandler) GetDefaultSmsConfig(ctx *gin.Context) {
	// 获取所有配置
	configs, err := h.smsService.GetAllSmsConfigs(context.Background())
	if err != nil {
		h.logError(ctx, "获取默认短信配置", "查询失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	// 查找默认配置（启用状态的配置）
	var defaultConfig *sms_domain.SmsConfig
	for _, config := range configs {
		if config.IsEnabled() {
			defaultConfig = config
			break
		}
	}

	if defaultConfig == nil {
		h.logError(ctx, "获取默认短信配置", "默认配置不存在", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("默认配置不存在"))
		return
	}

	h.logSuccess(ctx, "获取默认短信配置", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取默认短信配置成功", defaultConfig))
}

// GetSmsTemplates 获取短信模板
// 支持通过ID查询单个模板、通过标识查询模板或获取所有模板
func (h *SmsHandler) GetSmsTemplates(ctx *gin.Context) {
	id := ctx.Query("id")
	identifier := ctx.Query("identifier")

	// 如果提供了ID，获取单个模板
	if id != "" {
		template, err := h.smsService.GetSmsTemplate(context.Background(), id)
		if err != nil {
			h.logError(ctx, "获取短信模板", "查询失败", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
			return
		}

		if template == nil {
			h.logError(ctx, "获取短信模板", "模板不存在", nil)
			ctx.JSON(http.StatusOK, res.FailWithWarn("模板不存在"))
			return
		}

		h.logSuccess(ctx, "获取短信模板", "成功", map[string]interface{}{
			"template_id": id,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取短信模板成功", template))
		return
	}

	// 如果提供了标识，根据标识获取模板
	if identifier != "" {
		template, err := h.smsService.GetSmsTemplateByIdentifier(context.Background(), identifier)
		if err != nil {
			h.logError(ctx, "获取短信模板", "查询失败", err)
			ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
			return
		}

		if template == nil {
			h.logError(ctx, "获取短信模板", "模板不存在", nil)
			ctx.JSON(http.StatusOK, res.FailWithWarn("模板不存在"))
			return
		}

		h.logSuccess(ctx, "获取短信模板", "成功", map[string]interface{}{
			"template_identifier": identifier,
		})
		ctx.JSON(http.StatusOK, res.SuccessWithData("获取短信模板成功", template))
		return
	}

	// 否则获取所有模板
	templates, err := h.smsService.GetAllSmsTemplates(context.Background())
	if err != nil {
		h.logError(ctx, "获取所有短信模板", "查询失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "获取所有短信模板", "成功", nil)
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取所有短信模板成功", templates))
}

// UpdateSmsTemplateStatus 更新短信模板状态
func (h *SmsHandler) UpdateSmsTemplateStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Query("status")

	// 验证状态参数
	if status != "enable" && status != "disable" {
		h.logError(ctx, "更新短信模板状态", "状态参数无效", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("状态参数无效，应为'enable'或'disable'"))
		return
	}

	// 获取当前模板
	template, err := h.smsService.GetSmsTemplate(context.Background(), id)
	if err != nil {
		h.logError(ctx, "更新短信模板状态", "获取模板失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	if template == nil {
		h.logError(ctx, "更新短信模板状态", "模板不存在", nil)
		ctx.JSON(http.StatusOK, res.FailWithWarn("模板不存在"))
		return
	}

	// 根据状态更新模板
	if status == "enable" {
		template.Enable()
	} else {
		template.Disable()
	}

	// 更新模板
	if err := h.smsService.UpdateSmsTemplate(context.Background(), template); err != nil {
		h.logError(ctx, "更新短信模板状态", "更新模板失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新短信模板状态", "成功", map[string]interface{}{
		"template_id": id,
		"status":      status,
	})
	ctx.JSON(http.StatusOK, res.Success("状态更新成功"))
}

// CreateSmsTemplate 创建短信模板
func (h *SmsHandler) CreateSmsTemplate(ctx *gin.Context) {
	var req sms_domain.SmsTemplate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logError(ctx, "创建短信模板", "参数绑定失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	// 通过服务层创建模板
	if err := h.smsService.CreateSmsTemplate(context.Background(), &req); err != nil {
		h.logError(ctx, "创建短信模板", "保存失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "创建短信模板", "成功", map[string]interface{}{
		"template_id": req.SMB001,
	})
	ctx.JSON(http.StatusOK, res.SuccessWithData("创建短信模板成功", req))
}

// UpdateSmsTemplate 更新短信模板
func (h *SmsHandler) UpdateSmsTemplate(ctx *gin.Context) {
	id := ctx.Param("id")
	var req sms_domain.SmsTemplate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logError(ctx, "更新短信模板", "参数绑定失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	// 设置更新时间
	req.SMB008 = time.Now()

	req.SMB001 = id
	// 通过服务层更新模板
	if err := h.smsService.UpdateSmsTemplate(context.Background(), &req); err != nil {
		h.logError(ctx, "更新短信模板", "更新失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "更新短信模板", "成功", map[string]interface{}{
		"template_id": req.SMB001,
	})
	ctx.JSON(http.StatusOK, res.SuccessWithData("更新短信模板成功", req))
}

// DeleteSmsTemplate 删除短信模板
func (h *SmsHandler) DeleteSmsTemplate(ctx *gin.Context) {
	id := ctx.Param("id")

	// 通过服务层删除模板
	if err := h.smsService.DeleteSmsTemplate(context.Background(), id); err != nil {
		h.logError(ctx, "删除短信模板", "删除失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "删除短信模板", "成功", map[string]interface{}{
		"template_id": id,
	})
	ctx.JSON(http.StatusOK, res.Success("删除成功"))
}

// SearchSmsTemplates 搜索短信模板
func (h *SmsHandler) SearchSmsTemplates(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	// 简化实现，这里直接返回所有模板
	// 实际项目中可以根据keyword进行搜索
	templates, err := h.smsService.GetAllSmsTemplates(context.Background())
	if err != nil {
		h.logError(ctx, "搜索短信模板", "查询失败", err)
		ctx.JSON(http.StatusOK, res.FailWithWarn(err.Error()))
		return
	}

	h.logSuccess(ctx, "搜索短信模板", "成功", map[string]interface{}{
		"keyword": keyword,
	})
	ctx.JSON(http.StatusOK, res.SuccessWithData("获取成功", templates))
}

// 记录错误日志
func (h *SmsHandler) logError(ctx *gin.Context, module, content string, err error) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	logError := ""
	if err != nil {
		logError = err.Error()
	}
	h.logService.LogBusiness(ctx, log_domain.LogLevelError, "短信Web层", module, "", "", "",
		content, "失败", logError, clientIP, nil)
}

// 记录成功日志
func (h *SmsHandler) logSuccess(ctx *gin.Context, module, content string, extra map[string]interface{}) {
	clientIP, _ := middleware.GetRequestContextInfo(ctx)
	h.logService.LogBusiness(ctx, log_domain.LogLevelInfo, "短信Web层", module, "", "", "",
		content, "成功", "", clientIP, extra)
}
