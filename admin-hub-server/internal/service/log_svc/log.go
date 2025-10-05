package log_svc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/log_repo"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LogService 日志服务接口
type LogService interface {
	// LogCommon 记录常用日志
	// ctx: 上下文对象，用于传递请求范围的值
	// level: 日志级别（可选值：debug, info, warn, error）
	// message: 日志消息内容
	// fields: 额外的日志字段，用于记录更详细的信息
	// caller: 调用者信息，标识哪个服务或模块调用了日志记录
	// environment: 环境信息，如dev, test, prod等
	LogCommon(ctx context.Context, level log_domain.LogLevel, message string, fields map[string]interface{}, caller string, environment string) error

	// LogLogin 记录登录日志
	// ctx: 上下文对象
	// level: 日志级别（可选值：debug, info, warn, error）
	// message: 日志消息
	// userID: 用户ID
	// username: 用户名
	// ip: 用户IP地址
	// device: 用户设备信息
	// status: 登录状态（成功/失败）
	// reason: 状态原因描述
	// fields: 额外字段
	LogLogin(ctx context.Context, level log_domain.LogLevel, message string, userID, username, ip, device, status, reason string, fields map[string]interface{}) error

	// LogBusiness 记录业务日志
	// ctx: 上下文对象
	// level: 日志级别（可选值：debug, info, warn, error）
	// module: 业务模块名称
	// operation: 操作名称
	// businessID: 业务ID，标识具体的业务对象
	// operatorID: 操作者ID
	// operatorName: 操作者名称
	// content: 业务内容描述
	// status: 业务操作状态
	// reason: 状态原因描述
	// ip: 操作者IP地址
	// fields: 额外字段
	LogBusiness(ctx context.Context, level log_domain.LogLevel, module, operation, businessID, operatorID, operatorName, content, status, reason, ip string, fields map[string]interface{}) error

	// GetCommonLogs 获取常用日志列表
	// ctx: 上下文对象
	// page: 页码，从1开始
	// pageSize: 每页记录数
	// startTime: 开始时间，格式YYYY-MM-DD HH:MM:SS
	// endTime: 结束时间，格式YYYY-MM-DD HH:MM:SS
	GetCommonLogs(ctx context.Context, page, pageSize int, startTime, endTime string) ([]*log_domain.LL01, int64, error)

	// GetLoginLogs 获取登录日志列表
	// ctx: 上下文对象
	// page: 页码
	// pageSize: 每页记录数
	// startTime: 开始时间
	// endTime: 结束时间
	// userID: 用户ID，可选，用于筛选特定用户的登录日志
	GetLoginLogs(ctx context.Context, page, pageSize int, startTime, endTime string, userID string) ([]*log_domain.LL02, int64, error)

	// GetBusinessLogs 获取业务日志列表
	// ctx: 上下文对象
	// page: 页码
	// pageSize: 每页记录数
	// startTime: 开始时间
	// endTime: 结束时间
	// module: 业务模块，可选
	// operatorID: 操作者ID，可选
	GetBusinessLogs(ctx context.Context, page, pageSize int, startTime, endTime string, module, operatorID string) ([]*log_domain.LL03, int64, error)

	// DeleteLog 删除日志
	// ctx: 上下文对象
	// logType: 日志类型（可选值：common, business, login）
	// logID: 日志ID
	DeleteLog(ctx context.Context, logType log_domain.LogType, logID string) error
}

// logService 日志服务实现

type logService struct {
	repo   log_repo.LogRepository
	logger *zap.Logger
}

// NewLogService 创建日志服务实例
func NewLogService(repo log_repo.LogRepository, logger *zap.Logger) LogService {
	return &logService{
		repo:   repo,
		logger: logger,
	}
}

// LogCommon 记录常用日志
func (s *logService) LogCommon(ctx context.Context, level log_domain.LogLevel, message string, fields map[string]interface{}, caller string, environment string) error {
	// 生成日志ID
	logID := uuid.New().String()

	// 序列化字段
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	// 创建日志对象
	log := &log_domain.LL01{
		LLA001: logID,
		LLA002: level,
		LLA003: message,
		LLA004: string(fieldsJSON),
		LLA005: caller,
		LLA006: environment,
		LLA007: time.Now(),
	}

	// 异步保存到数据库
	go func() {
		if err := s.repo.SaveCommonLog(ctx, log); err != nil {
			s.logger.Error("保存常用日志失败", zap.Error(err), zap.String("logID", logID))
		}
	}()

	// 根据日志级别输出到控制台或文件
	switch level {
	case log_domain.LogLevelDebug:
		s.logger.Debug(message, zap.Any("fields", fields))
	case log_domain.LogLevelInfo:
		s.logger.Info(message, zap.Any("fields", fields))
	case log_domain.LogLevelWarn:
		s.logger.Warn(message, zap.Any("fields", fields))
	case log_domain.LogLevelError:
		s.logger.Error(message, zap.Any("fields", fields))
	}

	return nil
}

// LogLogin 记录登录日志
func (s *logService) LogLogin(ctx context.Context, level log_domain.LogLevel, message string, userID, username, ip, device, status, reason string, fields map[string]interface{}) error {
	// 生成日志ID
	logID := uuid.New().String()

	// 序列化字段
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	// 创建登录日志对象
	log := &log_domain.LL02{
		LLA001: logID,
		LLA002: level,
		LLA003: message,
		LLA004: userID,
		LLA005: username,
		LLA006: ip,
		LLA007: device,
		LLA008: status,
		LLA009: reason,
		LLA010: string(fieldsJSON),
		LLA011: time.Now(),
	}

	// 异步保存到数据库
	go func() {
		if err := s.repo.SaveLoginLog(ctx, log); err != nil {
			s.logger.Error("保存登录日志失败", zap.Error(err), zap.String("logID", logID))
		}
	}()

	// 根据日志级别输出到控制台或文件
	switch level {
	case log_domain.LogLevelDebug:
		s.logger.Debug(message, zap.String("userID", userID), zap.String("username", username), zap.String("ip", ip), zap.String("status", status))
	case log_domain.LogLevelInfo:
		s.logger.Info(message, zap.String("userID", userID), zap.String("username", username), zap.String("ip", ip), zap.String("status", status))
	case log_domain.LogLevelWarn:
		s.logger.Warn(message, zap.String("userID", userID), zap.String("username", username), zap.String("ip", ip), zap.String("status", status), zap.String("reason", reason))
	case log_domain.LogLevelError:
		s.logger.Error(message, zap.String("userID", userID), zap.String("username", username), zap.String("ip", ip), zap.String("status", status), zap.String("reason", reason))
	}

	return nil
}

// LogBusiness 记录业务日志
func (s *logService) LogBusiness(ctx context.Context, level log_domain.LogLevel, module, operation, businessID, operatorID, operatorName, content, status, reason, ip string, fields map[string]interface{}) error {
	// 使用content作为日志消息
	message := content

	// 生成日志ID
	logID := uuid.New().String()

	// 序列化字段
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	// 创建业务日志对象
	log := &log_domain.LL03{
		LLA001: logID,
		LLA002: level,
		LLA003: module,
		LLA004: operation,
		LLA005: businessID,
		LLA006: operatorID,
		LLA007: operatorName,
		LLA008: content,
		LLA009: status,
		LLA010: reason,
		LLA011: ip,
		LLA012: string(fieldsJSON),
		LLA013: time.Now(),
	}

	// 异步保存到数据库
	go func() {
		if err := s.repo.SaveBusinessLog(ctx, log); err != nil {
			s.logger.Error("保存业务日志失败", zap.Error(err), zap.String("logID", logID))
		}
	}()

	// 根据日志级别输出到控制台或文件
	switch level {
	case log_domain.LogLevelDebug:
		s.logger.Debug(message, zap.String("module", module), zap.String("operation", operation), zap.String("operatorID", operatorID), zap.String("status", status))
	case log_domain.LogLevelInfo:
		s.logger.Info(message, zap.String("module", module), zap.String("operation", operation), zap.String("operatorID", operatorID), zap.String("status", status))
	case log_domain.LogLevelWarn:
		s.logger.Warn(message, zap.String("module", module), zap.String("operation", operation), zap.String("operatorID", operatorID), zap.String("status", status), zap.String("reason", reason))
	case log_domain.LogLevelError:
		s.logger.Error(message, zap.String("module", module), zap.String("operation", operation), zap.String("operatorID", operatorID), zap.String("status", status), zap.String("reason", reason))
	}

	return nil
}

// GetCommonLogs 获取常用日志列表
func (s *logService) GetCommonLogs(ctx context.Context, page, pageSize int, startTime, endTime string) ([]*log_domain.LL01, int64, error) {
	// 调用Repository层获取数据
	// 这里需要转换dao模型到domain模型
	// 简化处理，实际实现中应该有完整的转换逻辑
	return nil, 0, nil
}

// GetLoginLogs 获取登录日志列表
func (s *logService) GetLoginLogs(ctx context.Context, page, pageSize int, startTime, endTime string, userID string) ([]*log_domain.LL02, int64, error) {
	// 调用Repository层获取数据
	// 这里需要转换dao模型到domain模型
	// 简化处理，实际实现中应该有完整的转换逻辑
	return nil, 0, nil
}

// GetBusinessLogs 获取业务日志列表
func (s *logService) GetBusinessLogs(ctx context.Context, page, pageSize int, startTime, endTime string, module, operatorID string) ([]*log_domain.LL03, int64, error) {
	// 调用Repository层获取数据
	// 这里需要转换dao模型到domain模型
	// 简化处理，实际实现中应该有完整的转换逻辑
	return nil, 0, nil
}

// DeleteLog 删除日志
func (s *logService) DeleteLog(ctx context.Context, logType log_domain.LogType, logID string) error {
	return s.repo.DeleteLog(ctx, logType, logID)
}
