package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

// Builder 定义了一个中间件构建器，用于配置和构建滑动窗口限流中间件
type Builder struct {
	prefix   string        // Redis 键的前缀
	cmd      redis.Cmdable // Redis 命令执行器
	interval time.Duration // 限流时间窗口
	// 阈值
	rate int // 限流阈值
}

//go:embed slide_window.lua
var luaScript string // 嵌入的 Lua 脚本

// NewBuilder 创建一个新的 Builder 实例
func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		cmd:      cmd,          // 设置 Redis 命令执行器
		prefix:   "ip-limiter", // 设置 Redis 键的前缀
		interval: interval,     // 设置限流时间窗口
		rate:     rate,         // 设置限流阈值
	}
}

// Prefix 设置 Redis 键的前缀
func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix // 设置 Redis 键的前缀
	return b
}

// Build 构建滑动窗口限流中间件
func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx) // 检查是否限流
		if err != nil {
			// 如果出错，返回 500 错误
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			// 如果限流，返回 429 错误
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next() // 继续处理请求
	}
}

// limit 检查是否限流
func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()) // 生成 Redis 键
	// 使用 Lua 脚本执行限流逻辑
	return b.cmd.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
