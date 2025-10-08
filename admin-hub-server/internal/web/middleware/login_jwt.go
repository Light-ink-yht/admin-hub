package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/service/user_svc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LoginJWTMiddlewareBuilder JWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string // 不需要登录校验的路径列表
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

// IgnorePaths 添加不需要登录校验的路径
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// Build 构建 JWT 登录校验中间件
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 我现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &user_svc.UserClaims{}
		// ParseWithClaims 里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return user_svc.JWTKey, nil
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid || claims.UserId == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()

		if claims.ExpiresAt.Sub(now) < time.Second*60 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
			tokenStr, err = token.SignedString(user_svc.JWTKey)
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
	}
}
