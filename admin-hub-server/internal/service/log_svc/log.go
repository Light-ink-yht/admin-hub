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
	// 记录常用日志
	LogCommon(ctx context.Context, level log_domain.LogLevel, message string, fields map[string]interface{}, caller string, environment string) error
	// 记录登录日志
	LogLogin(ctx context.Context, level log_domain.LogLevel, message string, userID, username, ip, device, status, reason string, fields map[string]interface{}) error
	// 记录业务日志
	LogBusiness(ctx context.Context, level log_domain.LogLevel, module, operation, businessID, operatorID, operatorName, content, status, reason, ip string, fields map[string]interface{}) error
	// 获取常用日志列表
	GetCommonLogs(ctx context.Context, page, pageSize int, startTime, endTime string) ([]*log_domain.LL01, int64, error)
	// 获取登录日志列表
	GetLoginLogs(ctx context.Context, page, pageSize int, startTime, endTime string, userID string) ([]*log_domain.LL02, int64, error)
	// 获取业务日志列表
	GetBusinessLogs(ctx context.Context, page, pageSize int, startTime, endTime string, module, operatorID string) ([]*log_domain.LL03, int64, error)
	// 删除日志
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
