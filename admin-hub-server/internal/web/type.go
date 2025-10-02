package web

import "github.com/gin-gonic/gin"

// Handler 定义了一个接口，用于注册路由
type Handler interface {
	// RegisterRoutes 注册路由
	RegisterRoutes(s *gin.RouterGroup)
}
