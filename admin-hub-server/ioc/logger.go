package ioc

import (
	"database/sql"
	pkgLogger "github.com/Light-ink-yht/admin-hub/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	Environment string `mapstructure:"environment"` // dev或prod
	Level       string `mapstructure:"level"`       // 日志级别
	FilePath    string `mapstructure:"file_path"`   // 日志文件路径
	ErrorPath   string `mapstructure:"error_path"`  // 错误日志路径
	UseDB       bool   `mapstructure:"use_db"`      // 是否使用数据库存储
}

// InitLogger 初始化日志系统并返回常用日志记录器
func InitLogger(db *sql.DB) *zap.Logger {
	// 初始化前创建一个临时logger用于记录初始化过程中的日志
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	// 从Viper中读取日志配置
	var logConfig LogConfig
	if err := viper.UnmarshalKey("log", &logConfig); err != nil {
		zap.L().Error("解析日志配置失败", zap.Error(err))
	}

	// 默认配置
	if logConfig.Environment == "" {
		logConfig.Environment = "dev"
	}
	if logConfig.Level == "" {
		logConfig.Level = "debug"
	}

	// 转换为logger包的配置格式
	newLogConfig := pkgLogger.LogConfig{
		Environment: logConfig.Environment,
		Level:       logConfig.Level,
		FilePath:    logConfig.FilePath,
		ErrorPath:   logConfig.ErrorPath,
		UseDB:       logConfig.UseDB,
	}

	// 初始化日志系统
	pkgLogger.Init(newLogConfig, db)
	// 获取常用日志记录器
	commonLogger := pkgLogger.GetCommonLogger()
	zap.ReplaceGlobals(commonLogger)

	// 如果使用数据库存储，确保日志表已创建
	if logConfig.UseDB && db != nil {
		if err := pkgLogger.EnsureLogTables(db); err != nil {
			commonLogger.Error("创建日志表失败", zap.Error(err))
		} else {
			commonLogger.Info("数据库日志存储已启用，日志表已确保存在")
		}
	}

	return commonLogger
}
