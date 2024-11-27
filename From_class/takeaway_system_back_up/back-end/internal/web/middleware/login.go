package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path ...string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path...)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		role := sess.Get("role") // 获取用户角色
		// 根据角色检查对应的ID
		switch role {
		case "customer":
			id := sess.Get("customerId")
			if id == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "customer not logged in"})
				return
			}
		case "employee":
			id := sess.Get("employeeId")
			if id == nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "employee not logged in"})
				return
			}
		default:
			// 如果没有角色或角色不是预期的值，则认为没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
