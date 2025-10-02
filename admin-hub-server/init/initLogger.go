package init

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig 日志配置结构体，与配置文件对应
type LogConfig struct {
	Environment string `mapstructure:"environment"` // dev或prod
	Level       string `mapstructure:"level"`       // 日志级别
	FilePath    string `mapstructure:"file_path"`   // 日志文件路径
	ErrorPath   string `mapstructure:"error_path"`  // 错误日志路径
}

func Logger() {
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

	// 解析日志级别
	level := zapcore.InfoLevel
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	var logger *zap.Logger
	var err error

	if logConfig.Environment == "dev" {
		// 开发环境配置
		cfg := zap.NewDevelopmentConfig()
		cfg.Level.SetLevel(level)
		logger, err = cfg.Build(zap.AddCaller())
	} else {
		// 生产环境配置
		cfg := zap.NewProductionConfig()
		cfg.Level.SetLevel(level)

		// 设置日志输出路径
		outputPaths := []string{"stderr"}
		if logConfig.FilePath != "" {
			outputPaths = append(outputPaths, logConfig.FilePath)
		}
		cfg.OutputPaths = outputPaths

		// 设置错误日志路径
		errorOutputPaths := []string{"stderr"}
		if logConfig.ErrorPath != "" {
			errorOutputPaths = append(errorOutputPaths, logConfig.ErrorPath)
		}
		cfg.ErrorOutputPaths = errorOutputPaths

		// 自定义时间格式
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, err = cfg.Build(zap.AddCaller())
	}

	if err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}

	// 替换全局日志实例
	zap.ReplaceGlobals(logger)
	zap.L().Info("日志初始化完成", zap.String("环境", logConfig.Environment), zap.String("级别", logConfig.Level))
}
