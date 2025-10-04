package appinit

import (
	"database/sql"
	"fmt"

	"github.com/Light-ink-yht/admin-hub/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LogConfig 日志配置结构体，与配置文件对应
type LogConfig struct {
	Environment string `mapstructure:"environment"` // dev或prod
	Level       string `mapstructure:"level"`       // 日志级别
	FilePath    string `mapstructure:"file_path"`   // 日志文件路径
	ErrorPath   string `mapstructure:"error_path"`  // 错误日志路径
	UseDB       bool   `mapstructure:"use_db"`      // 是否使用数据库存储
}

// Logger 初始化日志系统
func Logger(db *sql.DB) {
	// 从Viper中读取日志配置
	var logConfig LogConfig
	if err := viper.UnmarshalKey("log", &logConfig); err != nil {
		panic(fmt.Sprintf("解析日志配置失败: %v", err))
	}

	// 默认配置
	if logConfig.Environment == "" {
		logConfig.Environment = "dev"
	}
	if logConfig.Level == "" {
		logConfig.Level = "debug"
	}

	// 转换为我们新日志系统的配置格式
	newLogConfig := logger.LogConfig{
		Environment: logConfig.Environment,
		Level:       logConfig.Level,
		FilePath:    logConfig.FilePath,
		ErrorPath:   logConfig.ErrorPath,
		UseDB:       logConfig.UseDB,
	}

	// 初始化新的日志系统
	logger.Init(newLogConfig, db)

	// 测试日志
	commonLogger := logger.GetCommonLogger()
	commonLogger.Info("常用日志系统初始化完成")

	businessLogger := logger.GetBusinessLogger()
	businessLogger.Info("业务日志系统初始化完成")

	loginLogger := logger.GetLoginLogger()
	loginLogger.Info("登录日志系统初始化完成")

	// 如果使用数据库存储，确保日志表已创建
	if logConfig.UseDB && db != nil {
		// 确保日志表存在
		if err := logger.EnsureLogTables(db); err != nil {
			zap.L().Error("创建日志表失败", zap.Error(err))
		} else {
			zap.L().Info("数据库日志存储已启用，日志表已确保存在")
		}
	}
}
