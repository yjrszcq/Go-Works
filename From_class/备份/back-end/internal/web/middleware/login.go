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
		/*
			或直接这样
			if ctx.Request.URL.Path == "/users/login" ||
					ctx.Request.URL.Path == "/users/signup" {
					return
				}
		*/
		sess := sessions.Default(ctx)
		/*
			if sess == nil {
				// 没有登录
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			// 不可能为 nil，因为上面已经有 session 插件了
		*/
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
