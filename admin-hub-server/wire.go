//go:build wireinject

package main

import (
	"github.com/Light-ink-yht/admin-hub/internal/repository/cache/code_cache"
	"github.com/Light-ink-yht/admin-hub/internal/repository/code_repo"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
	"github.com/Light-ink-yht/admin-hub/internal/repository/email_repo"
	"github.com/Light-ink-yht/admin-hub/internal/repository/log_repo"
	"github.com/Light-ink-yht/admin-hub/internal/service/code_svc"
	"github.com/Light-ink-yht/admin-hub/internal/service/log_svc"
	"github.com/Light-ink-yht/admin-hub/internal/web/email_web"
	"github.com/Light-ink-yht/admin-hub/internal/web/user_web"
	"github.com/Light-ink-yht/admin-hub/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// InitWebServer 初始化 Web 服务器
func InitWebServer() *gin.Engine {
	// 使用 Wire 依赖注入框架来构建依赖关系
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitLogger,

		// 初始化 DAO
		dao.NewCodeDAO,
		dao.NewEmailConfigDAO,
		dao.NewEmailTemplateDAO,

		// 初始化 cache
		code_cache.NewCodeCache,

		// 初始化 Repository
		code_repo.NewCodeRepository,
		email_repo.NewEmailConfigRepository,
		email_repo.NewEmailTemplateRepository,

		// 初始化 LogService 相关依赖
		log_repo.NewLogRepository,
		log_svc.NewLogService,

		// 初始化 Service
		code_svc.NewEmailService,

		// 初始化 Handler
		user_web.NewUserHandler,   // 会自动注入 codeSvc, logService, logger
		email_web.NewEmailHandler, // 会自动注入 emailConfigRepo, emailTemplateRepo, logService

		// 初始化 Web 服务器
		ioc.InitWebServer,
	)
	// 返回一个新的 Gin 引擎实例
	return new(gin.Engine)
}
