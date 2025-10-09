package logger

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// 全局日志记录器
	globalLogger *zap.Logger
	mu           sync.RWMutex

	// 特定类型的日志记录器
	commonLogger   *zap.Logger
	businessLogger *zap.Logger
	loginLogger    *zap.Logger
)

// LogType 日志类型枚举
const (
	LogTypeCommon   = "ll01" // 常用日志
	LogTypeBusiness = "ll02" // 业务日志
	LogTypeLogin    = "ll03" // 登录日志
)

// LogConfig 日志配置
type LogConfig struct {
	Environment string `mapstructure:"environment"`
	Level       string `mapstructure:"level"`
	FilePath    string `mapstructure:"file_path"`
	ErrorPath   string `mapstructure:"error_path"`
	UseDB       bool   `mapstructure:"use_db"`
	MaxFileSize int    `mapstructure:"max_file_size"` // MB
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"` // days
	Compress    bool   `mapstructure:"compress"`
}

// Init 初始化日志系统
func Init(cfg LogConfig, db *sql.DB) {
	// 解析日志级别
	level := ConvertStringToLevel(cfg.Level)

	var logger *zap.Logger
	var err error
	var cores []zapcore.Core

	// 创建编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	if cfg.Environment == "dev" {
		// 开发环境配置
		devCfg := zap.NewDevelopmentConfig()
		devCfg.Level.SetLevel(level)
		logger, err = devCfg.Build(zap.AddCaller())
	} else {
		// 生产环境配置
		// 控制台输出
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, consoleCore)

		// 文件输出
		if cfg.FilePath != "" {
			fileWriter, err, _ := zap.Open(cfg.FilePath)
			if err != nil {
				panic(fmt.Sprintf("打开日志文件失败: %v", err))
			}
			fileCore := zapcore.NewCore(encoder, fileWriter, level)
			cores = append(cores, fileCore)

			// 错误日志单独输出
			if cfg.ErrorPath != "" {
				errorWriter, err, _ := zap.Open(cfg.ErrorPath)
				if err != nil {
					panic(fmt.Sprintf("打开错误日志文件失败: %v", err))
				}
				errorCore := zapcore.NewCore(encoder, errorWriter, zapcore.ErrorLevel)
				cores = append(cores, errorCore)
			}
		}

		// 数据库输出
		if cfg.UseDB && db != nil {
			dbCore := NewDBCore(db, level, encoder)
			cores = append(cores, dbCore)
		}

		// 合并多个核心
		core := zapcore.NewTee(cores...)

		// 创建logger
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	if err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}

	mu.Lock()
	globalLogger = logger
	mu.Unlock()

	// 初始化具体类型的logger
	commonLogger = logger.With(zap.String("log_type", LogTypeCommon))
	businessLogger = logger.With(zap.String("log_type", LogTypeBusiness))
	loginLogger = logger.With(zap.String("log_type", LogTypeLogin))

	logger.Info("日志初始化完成", zap.String("环境", cfg.Environment), zap.String("级别", cfg.Level))
}

// GetLogger 获取指定类型的日志记录器
func GetLogger(logType string) *zap.Logger {
	mu.RLock()
	defer mu.RUnlock()

	if globalLogger == nil {
		panic("日志系统尚未初始化")
	}

	// 根据日志类型返回对应的logger
	switch logType {
	case LogTypeCommon:
		return commonLogger
	case LogTypeBusiness:
		return businessLogger
	case LogTypeLogin:
		return loginLogger
	default:
		return globalLogger.With(zap.String("log_type", logType))
	}
}

// GetCommonLogger 获取常用日志记录器
func GetCommonLogger() *zap.Logger {
	return GetLogger(LogTypeCommon)
}

// GetBusinessLogger 获取业务日志记录器
func GetBusinessLogger() *zap.Logger {
	return GetLogger(LogTypeBusiness)
}

// GetLoginLogger 获取登录日志记录器
func GetLoginLogger() *zap.Logger {
	return GetLogger(LogTypeLogin)
}

// DBCore 实现 zapcore.Core 接口，用于将日志写入数据库
type DBCore struct {
	db           *sql.DB
	levelEnabler zapcore.LevelEnabler
	encoder      zapcore.Encoder
	context      context.Context
	mutex        sync.Mutex
}

// Enabled 实现 zapcore.Core 接口，检查日志级别是否应该被记录
func (c *DBCore) Enabled(level zapcore.Level) bool {
	return c.levelEnabler.Enabled(level)
}

// NewDBCore 创建一个新的 DBCore 实例
func NewDBCore(db *sql.DB, level zapcore.LevelEnabler, encoder zapcore.Encoder) *DBCore {
	return &DBCore{
		db:           db,
		levelEnabler: level,
		encoder:      encoder,
		context:      context.Background(),
	}
}

// With 实现 zapcore.Core 接口
func (c *DBCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.encoder, fields)
	return clone
}

// Check 实现 zapcore.Core 接口
func (c *DBCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.levelEnabler.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

// Write 实现 zapcore.Core 接口，将日志写入数据库
func (c *DBCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 构建字段map
	fieldsMap := MapFromFields(fields)

	// 确定日志类型
	logType := LogTypeCommon
	for _, field := range fields {
		if field.Key == "log_type" {
			if logTypeStr, ok := field.Interface.(string); ok {
				switch logTypeStr {
				case LogTypeCommon, LogTypeBusiness, LogTypeLogin:
					logType = logTypeStr
				}
			}
			break
		}
	}

	// 使用带重试的方式写入日志
	return LogWithRetry(zap.L(), func() error {
		// 根据日志类型选择表和参数
		switch logType {
		case LogTypeCommon:
			return c.writeCommonLog(entry, fieldsMap)
		case LogTypeBusiness:
			return c.writeBusinessLog(entry, fieldsMap)
		case LogTypeLogin:
			return c.writeLoginLog(entry, fieldsMap)
		default:
			return fmt.Errorf("未知的日志类型: %s", logType)
		}
	}, 3)
}

// writeCommonLog 写入常用日志
func (c *DBCore) writeCommonLog(entry zapcore.Entry, fieldsMap map[string]interface{}) error {
	// 生成日志ID
	logID := NewLoggerID()

	// 日志级别字符串
	levelStr := ConvertLevelToString(entry.Level)

	// 调用位置
	caller := FormatCaller(entry.Caller)

	// 字段JSON
	fieldsJSON := JSONStringFromMap(fieldsMap)

	// 常用日志插入
	_, err := c.db.ExecContext(c.context,
		"INSERT INTO ll01 (lla001, lla002, lla003, lla004, lla005, lla006, lla007) VALUES (?, ?, ?, ?, ?, ?, ?)",
		logID,
		levelStr,
		entry.Message,
		fieldsJSON,
		caller,
		"", // 环境标识
		entry.Time,
	)
	return err
}

// writeLoginLog 写入登录日志
func (c *DBCore) writeLoginLog(entry zapcore.Entry, fieldsMap map[string]interface{}) error {
	// 生成日志ID
	logID := NewLoggerID()

	// 日志级别字符串
	levelStr := ConvertLevelToString(entry.Level)

	// 从字段中提取登录相关信息
	userID := ""
	username := ""
	loginIP := ""
	device := ""
	status := "成功"
	failReason := ""

	if val, ok := fieldsMap["user_id"]; ok {
		userID = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["username"]; ok {
		username = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["ip"]; ok {
		loginIP = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["device"]; ok {
		device = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["status"]; ok {
		status = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["reason"]; ok {
		failReason = fmt.Sprintf("%v", val)
	}

	// 字段JSON
	fieldsJSON := JSONStringFromMap(fieldsMap)

	// 登录日志插入
	_, err := c.db.ExecContext(c.context,
		"INSERT INTO ll02 (lla001, lla002, lla003, lla004, lla005, lla006, lla007, lla008, lla009, lla010, lla011) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		logID,
		levelStr,
		entry.Message,
		userID,
		username,
		loginIP,
		device,
		status,
		failReason,
		fieldsJSON,
		entry.Time,
	)
	return err
}

// writeBusinessLog 写入业务日志
func (c *DBCore) writeBusinessLog(entry zapcore.Entry, fieldsMap map[string]interface{}) error {
	// 生成日志ID
	logID := NewLoggerID()

	// 日志级别字符串
	levelStr := ConvertLevelToString(entry.Level)

	// 从字段中提取业务相关信息
	module := ""
	opType := ""
	businessID := ""
	operatorID := ""
	operatorName := ""
	content := ""
	status := "成功"
	failReason := ""
	opIP := ""

	if val, ok := fieldsMap["module"]; ok {
		module = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["type"]; ok {
		opType = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["business_id"]; ok {
		businessID = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["operator_id"]; ok {
		operatorID = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["operator_name"]; ok {
		operatorName = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["content"]; ok {
		content = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["status"]; ok {
		status = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["reason"]; ok {
		failReason = fmt.Sprintf("%v", val)
	}
	if val, ok := fieldsMap["ip"]; ok {
		opIP = fmt.Sprintf("%v", val)
	}

	// 字段JSON
	fieldsJSON := JSONStringFromMap(fieldsMap)

	// 业务日志插入
	_, err := c.db.ExecContext(c.context,
		"INSERT INTO ll03 (lla001, lla002, lla003, lla004, lla005, lla006, lla007, lla008, lla009, lla010, lla011, lla012, lla013) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		logID,
		levelStr,
		module,
		opType,
		businessID,
		operatorID,
		operatorName,
		content,
		status,
		failReason,
		opIP,
		fieldsJSON,
		entry.Time,
	)
	return err
}

// Sync 实现 zapcore.Core 接口
func (c *DBCore) Sync() error {
	return nil
}

// 克隆 DBCore
func (c *DBCore) clone() *DBCore {
	return &DBCore{
		db:           c.db,
		levelEnabler: c.levelEnabler,
		encoder:      c.encoder.Clone(),
		context:      c.context,
	}
}

// 添加字段到编码器
func addFields(enc zapcore.ObjectEncoder, fields []zap.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
