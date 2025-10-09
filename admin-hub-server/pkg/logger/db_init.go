package logger

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// InitDBLogger 初始化数据库日志功能
func InitDBLogger(db *sql.DB, logLevel string) (*zap.Logger, error) {
	// 解析日志级别
	level := zapcore.InfoLevel
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// 创建编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建数据库日志核心
	dbCore := NewDBCore(db, level, encoder)

	// 创建多核心日志器，同时输出到控制台和数据库
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), level)

	// 使用Tee组合多个核心
	core := zapcore.NewTee(consoleCore, dbCore)

	// 创建logger
	logger := zap.New(core, zap.AddCaller())

	// 确保日志表存在
	err := ensureLogTables(db)
	if err != nil {
		return nil, fmt.Errorf("确保日志表存在失败: %w", err)
	}

	return logger, nil
}

// EnsureLogTables 确保日志表存在（公开函数）
func EnsureLogTables(db *sql.DB) error {
	return ensureLogTables(db)
}

// ensureLogTables 确保日志表存在（内部函数）
func ensureLogTables(db *sql.DB) error {
	// 创建常用日志表
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS ll01 (
            lan001 BIGINT AUTO_INCREMENT PRIMARY KEY,
            lan002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            lan003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            lan004 DATETIME,
            lan005 VARCHAR(1) NOT NULL DEFAULT '1',
            lan006 BIGINT NOT NULL DEFAULT 0,
            lan007 BIGINT NOT NULL DEFAULT 0,
            lla001 VARCHAR(50) NOT NULL COMMENT '日志ID',
            lla002 VARCHAR(10) NOT NULL COMMENT '日志级别',
            lla003 VARCHAR(255) NOT NULL COMMENT '日志消息',
            lla004 TEXT COMMENT '日志字段（JSON格式）',
            lla005 VARCHAR(255) COMMENT '调用位置',
            lla006 VARCHAR(20) COMMENT '环境标识',
            lla007 DATETIME NOT NULL COMMENT '日志时间',
            INDEX idx_lla007 (lla007),
            INDEX idx_lla001 (lla001),
            INDEX idx_lla002 (lla002)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='常用日志表';
    `)
	if err != nil {
		return fmt.Errorf("创建常用日志表失败: %w", err)
	}

	// 创建登录日志表
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ll02 (
            lan001 BIGINT AUTO_INCREMENT PRIMARY KEY,
            lan002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            lan003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            lan004 DATETIME,
            lan005 VARCHAR(1) NOT NULL DEFAULT '1',
            lan006 BIGINT NOT NULL DEFAULT 0,
            lan007 BIGINT NOT NULL DEFAULT 0,
            lla001 VARCHAR(50) NOT NULL COMMENT '日志ID',
            lla002 VARCHAR(10) NOT NULL COMMENT '日志级别',
            lla003 VARCHAR(255) NOT NULL COMMENT '日志消息',
            lla004 VARCHAR(50) COMMENT '用户ID',
            lla005 VARCHAR(50) COMMENT '用户名',
            lla006 VARCHAR(50) COMMENT '登录IP',
            lla007 VARCHAR(255) COMMENT '登录设备',
            lla008 VARCHAR(10) COMMENT '登录状态（成功/失败）',
            lla009 VARCHAR(255) COMMENT '失败原因',
            lla010 TEXT COMMENT '日志字段（JSON格式）',
            lla011 DATETIME NOT NULL COMMENT '登录时间',
            INDEX idx_lla011 (lla011),
            INDEX idx_lla001 (lla001),
            INDEX idx_lla004 (lla004),
            INDEX idx_lla006 (lla006)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';
    `)
	if err != nil {
		return fmt.Errorf("创建登录日志表失败: %w", err)
	}

	// 创建业务日志表
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ll03 (
            lan001 BIGINT AUTO_INCREMENT PRIMARY KEY,
            lan002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            lan003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            lan004 DATETIME,
            lan005 VARCHAR(1) NOT NULL DEFAULT '1',
            lan006 BIGINT NOT NULL DEFAULT 0,
            lan007 BIGINT NOT NULL DEFAULT 0,
            lla001 VARCHAR(50) NOT NULL COMMENT '日志ID',
            lla002 VARCHAR(10) NOT NULL COMMENT '日志级别',
            lla003 VARCHAR(50) COMMENT '业务模块',
            lla004 VARCHAR(50) COMMENT '操作类型',
            lla005 VARCHAR(50) COMMENT '业务ID',
            lla006 VARCHAR(50) COMMENT '操作人ID',
            lla007 VARCHAR(50) COMMENT '操作人姓名',
            lla008 VARCHAR(255) COMMENT '操作内容',
            lla009 VARCHAR(10) COMMENT '操作状态（成功/失败）',
            lla010 VARCHAR(255) COMMENT '失败原因',
            lla011 VARCHAR(50) COMMENT '操作IP',
            lla012 TEXT COMMENT '日志字段（JSON格式）',
            lla013 DATETIME NOT NULL COMMENT '操作时间',
            INDEX idx_lla013 (lla013),
            INDEX idx_lla001 (lla001),
            INDEX idx_lla003 (lla003),
            INDEX idx_lla006 (lla006)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='业务日志表';
    `)
	if err != nil {
		return fmt.Errorf("创建业务日志表失败: %w", err)
	}

	return nil
}
