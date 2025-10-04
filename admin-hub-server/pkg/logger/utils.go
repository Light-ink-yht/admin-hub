package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLoggerID 生成新的日志ID
func NewLoggerID() string {
	return uuid.New().String()
}

// FieldFromMap 将map转换为zap字段
func FieldFromMap(fields map[string]interface{}) []zap.Field {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// MapFromFields 将zap字段转换为map
func MapFromFields(fields []zap.Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Interface
	}
	return result
}

// JSONStringFromMap 将map转换为JSON字符串
func JSONStringFromMap(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}

	return string(data)
}

// ConvertLevelToString 将zap日志级别转换为字符串
func ConvertLevelToString(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return "debug"
	case zapcore.InfoLevel:
		return "info"
	case zapcore.WarnLevel:
		return "warn"
	case zapcore.ErrorLevel:
		return "error"
	case zapcore.DPanicLevel:
		return "dpanic"
	case zapcore.PanicLevel:
		return "panic"
	case zapcore.FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

// ConvertStringToLevel 将字符串转换为zap日志级别
func ConvertStringToLevel(levelStr string) zapcore.Level {
	var level zapcore.Level
	err := level.Set(levelStr)
	if err != nil {
		return zapcore.InfoLevel // 默认使用info级别
	}
	return level
}

// FormatCaller 格式化调用位置
func FormatCaller(caller zapcore.EntryCaller) string {
	if !caller.Defined {
		return "undefined"
	}
	return fmt.Sprintf("%s:%d", caller.TrimmedPath(), caller.Line)
}

// NewContext 创建包含追踪信息的上下文
func NewContext(requestID string, userID string) map[string]interface{} {
	return map[string]interface{}{
		"request_id": requestID,
		"user_id":    userID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
}

// WithTraceContext 添加追踪上下文到日志字段
func WithTraceContext(logger *zap.Logger, requestID string, userID string) *zap.Logger {
	return logger.With(
		zap.String("request_id", requestID),
		zap.String("user_id", userID),
	)
}

// LogWithRetry 带重试的日志记录函数
func LogWithRetry(logger *zap.Logger, fn func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
	}
	logger.Error("日志记录失败", zap.Error(err), zap.Int("max_retries", maxRetries))
	return err
}

// GetCurrentTime 获取当前时间（用于统一时间格式）
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串（用于统一时间格式）
func GetCurrentTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// NewLogIDWithPrefix 生成带前缀的日志ID
func NewLogIDWithPrefix(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, uuid.New().String())
}

// IsLogLevelEnabled 检查指定日志级别是否启用
func IsLogLevelEnabled(logger *zap.Logger, level zapcore.Level) bool {
	return logger.Core().Enabled(level)
}
