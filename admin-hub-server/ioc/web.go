package ioc

import (
	"github.com/Light-ink-yht/admin-hub/internal/web/user_web"
	"github.com/Light-ink-yht/admin-hub/pkg/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// InitWebServer 初始化 Web 服务器
func InitWebServer(mdls []gin.HandlerFunc, userHdl *user_web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	r := server.Group("/api")
	userHdl.RegisterRoutes(r)
	return server
}

// InitMiddlewares 初始化中间件
func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 跨域资源共享中间件
		corsHdl(),
		// 基于 Redis 的速率限制中间件
		ratelimit.NewBuilder(redisClient, time.Minute, 100).Build(),
	}
}

// corsHdl 处理跨域资源共享
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的源，这里注释掉了，实际使用时需要根据需求配置
		//AllowOrigins: []string{"*"},
		// 允许的方法，这里注释掉了，实际使用时需要根据需求配置
		//AllowMethods: []string{"POST", "GET"},
		// 允许的请求头
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 暴露的响应头，前端需要这个才能获取到自定义的响应头
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许携带凭证（如 cookie）
		AllowCredentials: true,
		// 自定义允许的源的函数
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 允许本地开发环境
				return true
			}
			// 允许特定域名
			return strings.Contains(origin, "yourcompany.com")
		},
		// 预检请求的有效期
		MaxAge: 12 * time.Hour,
	})
}
