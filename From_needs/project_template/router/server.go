package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"project/lib"
	"project/middleware"
	"strings"
	"time"
)

func InitServer(config *lib.Config) *gin.Engine {
	server := gin.Default()
	// 解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"}, // 可以不写，不写就是所有都支持
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // 是否允许带cookies之类的东西
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") {
				// 开发环境
				return true
			}
			if strings.Contains(origin, "http://127.0.0.1") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, config.AllowHost)
		},
		MaxAge: 12 * time.Second,
	}))
	// 登录校验
	store := cookie.NewStore([]byte("secret"))  // 存 sessions 数据的地方
	server.Use(sessions.Sessions("sid", store)) // cookie 的名字和值(store)
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/user/signup", "/user/login").
		IgnorePaths("/index").
		Build())

	return server
}
