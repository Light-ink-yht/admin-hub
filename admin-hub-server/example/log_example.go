package example

import (
	"fmt"

	"github.com/Light-ink-yht/admin-hub/pkg/logger"
	"go.uber.org/zap"
)

// LoggingExample 演示如何使用分类日志系统
func LoggingExample() {
	// 获取不同类型的日志记录器
	commonLogger := logger.GetCommonLogger()
	businessLogger := logger.GetBusinessLogger()
	loginLogger := logger.GetLoginLogger()

	// 记录常用日志
	commonLogger.Info("系统启动完成",
		zap.String("component", "system"),
		zap.String("version", "1.0.0"))

	commonLogger.Debug("处理请求",
		zap.String("path", "/api/users"),
		zap.String("method", "GET"))

	commonLogger.Error("配置加载失败",
		zap.String("file", "config.yaml"),
		zap.Error(fmt.Errorf("文件不存在")))

	// 记录业务日志
	businessLogger.Info("用户创建成功",
		zap.String("module", "user"),
		zap.String("type", "create"),
		zap.String("business_id", "1001"),
		zap.String("operator_id", "admin"),
		zap.String("operator_name", "管理员"),
		zap.String("content", "创建新用户"),
		zap.String("status", "成功"),
		zap.String("ip", "192.168.1.100"))

	businessLogger.Warn("业务操作失败",
		zap.String("module", "order"),
		zap.String("type", "pay"),
		zap.String("business_id", "2001"),
		zap.String("operator_id", "user1"),
		zap.String("operator_name", "用户1"),
		zap.String("content", "订单支付"),
		zap.String("status", "失败"),
		zap.String("reason", "余额不足"),
		zap.String("ip", "192.168.1.101"))

	// 记录登录日志
	loginLogger.Info("用户登录成功",
		zap.String("user_id", "admin"),
		zap.String("username", "管理员"),
		zap.String("ip", "192.168.1.200"),
		zap.String("device", "Chrome/90.0.4430.212"),
		zap.String("status", "成功"))

	loginLogger.Error("用户登录失败",
		zap.String("user_id", "unknown"),
		zap.String("username", "test"),
		zap.String("ip", "192.168.1.201"),
		zap.String("device", "Mozilla/5.0"),
		zap.String("status", "失败"),
		zap.String("reason", "密码错误"))

	// 使用工具函数
	ctx := logger.NewContext("req-123456", "user-789")
	fields := logger.FieldFromMap(ctx)
	commonLogger.Info("带上下文的日志", fields...)

	// 检查日志级别是否启用
	if logger.IsLogLevelEnabled(commonLogger, zap.DebugLevel) {
		commonLogger.Debug("此调试日志仅在调试级别启用时显示")
	}
}
