package log_domain

import (
	"time"

	"github.com/Light-ink-yht/admin-hub/pkg/model"
)

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogType 日志类型
type LogType string

const (
	// LogTypeCommon 常用日志
	LogTypeCommon LogType = "common"
	// LogTypeBusiness 业务日志
	LogTypeBusiness LogType = "business"
	// LogTypeLogin 登录日志
	LogTypeLogin LogType = "login"
)

// LL01 常用日志模型
type LL01 struct {
	model.LN01
	LLA001 string    // 日志ID
	LLA002 LogLevel  // 日志级别
	LLA003 string    // 日志消息
	LLA004 string    // 日志字段（JSON格式）
	LLA005 string    // 调用位置
	LLA006 string    // 环境标识
	LLA007 time.Time // 日志时间
}

// LL02 登录日志模型
type LL02 struct {
	model.LN01
	LLA001 string    // 日志ID
	LLA002 LogLevel  // 日志级别
	LLA003 string    // 日志消息
	LLA004 string    // 用户ID
	LLA005 string    // 用户名
	LLA006 string    // 登录IP
	LLA007 string    // 登录设备
	LLA008 string    // 登录状态（成功/失败）
	LLA009 string    // 失败原因
	LLA010 string    // 日志字段（JSON格式）
	LLA011 time.Time // 登录时间
}

// LL03 业务日志模型
type LL03 struct {
	model.LN01
	LLA001 string    // 日志ID
	LLA002 LogLevel  // 日志级别
	LLA003 string    // 业务模块
	LLA004 string    // 操作类型
	LLA005 string    // 业务ID
	LLA006 string    // 操作人ID
	LLA007 string    // 操作人姓名
	LLA008 string    // 操作内容
	LLA009 string    // 操作状态（成功/失败）
	LLA010 string    // 失败原因
	LLA011 string    // 操作IP
	LLA012 string    // 日志字段（JSON格式）
	LLA013 time.Time // 操作时间
}
