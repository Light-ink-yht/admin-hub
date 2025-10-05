package middleware

import (
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLogMiddleware 创建日志中间件
func NewLogMiddleware(logService log_svc.LogService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		// 获取请求信息
		clientIP := c.ClientIP()
		reqURL := c.Request.URL.Path
		reqMethod := c.Request.Method
		//userAgent := c.Request.UserAgent()
		//referer := c.Request.Referer()

		// 构建请求字段
		//reqFields := map[string]interface{}{
		//	"ip":         clientIP,
		//	"url":        reqURL,
		//	"method":     reqMethod,
		//	"user_agent": userAgent,
		//	"referer":    referer,
		//}

		// 记录请求接收日志
		logger.Info("收到HTTP请求", zap.String("url", reqURL), zap.String("method", reqMethod), zap.String("ip", clientIP))

		// 设置一些上下文信息，供后续处理使用
		c.Set("start_time", startTime)
		c.Set("client_ip", clientIP)

		// 处理请求
		c.Next()

		// 计算处理时间
		handleTime := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 构建响应字段
		//respFields := map[string]interface{}{
		//	"ip":          clientIP,
		//	"url":         reqURL,
		//	"method":      reqMethod,
		//	"status_code": statusCode,
		//	"duration_ms": handleTime.Milliseconds(),
		//}

		// 根据状态码确定日志级别
		logLevel := log_domain.LogLevelInfo
		if statusCode >= 400 && statusCode < 500 {
			logLevel = log_domain.LogLevelWarn
		} else if statusCode >= 500 {
			logLevel = log_domain.LogLevelError
		}

		// 记录响应日志
		message := "HTTP请求处理完成"
		if logLevel == log_domain.LogLevelWarn {
			message = "HTTP请求处理警告"
		} else if logLevel == log_domain.LogLevelError {
			message = "HTTP请求处理失败"
		}

		logger.Info(message, zap.Int("status_code", statusCode), zap.Duration("耗时", handleTime), zap.String("ip", clientIP), zap.String("url", reqURL))
	}
}

// GetRequestContextInfo 从上下文中获取请求相关信息，供业务层使用
func GetRequestContextInfo(c *gin.Context) (string, time.Duration) {
	// 获取客户端IP
	clientIP, exists := c.Get("client_ip")
	if !exists {
		clientIP = c.ClientIP()
	} else if ipStr, ok := clientIP.(string); ok {
		clientIP = ipStr
	} else {
		clientIP = c.ClientIP()
	}

	// 计算处理时间
	startTime, exists := c.Get("start_time")
	handleTime := time.Since(startTime.(time.Time))
	if !exists {
		handleTime = 0
	} else if startTimeValue, ok := startTime.(time.Time); ok {
		handleTime = time.Since(startTimeValue)
	}

	return clientIP.(string), handleTime
}
