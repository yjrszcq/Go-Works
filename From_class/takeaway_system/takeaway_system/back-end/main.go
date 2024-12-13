package main

import (
	"back-end/internal/repository/dao"
	"back-end/internal/web/init_web"
	"back-end/internal/web/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	cfg := init_web.InitConfig()
	db := initDB(cfg)
	server := initWebServer(cfg)
	init_web.InitWeb(db, server, cfg)
	server.Run(":" + cfg.ServerPort)
}

func initDB(config *init_web.Config) *gorm.DB {
	//配置MySQL连接参数
	username := config.DbUsername
	password := config.DbPassword
	host := config.DbHost
	port := config.DbPort
	Dbname := config.DbName
	timeout := config.DbTimeout
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	if config.DbInit {
		// 自动创建表
		err = dao.InitTable(db)
		if err != nil {
			panic(err)
		}
	}
	return db
}

func initWebServer(config *init_web.Config) *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{ // 解决跨域问题
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
	store := cookie.NewStore([]byte("secret"))              // 存 sessions 数据的地方
	server.Use(sessions.Sessions(config.CookieName, store)) // cookie 的名字和值(store)
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/customer/signup", "/customer/login").
		IgnorePaths("/employee/signup", "/employee/login").
		IgnorePaths("/admin/login", "/admin/logout").
		IgnorePaths("/dish/list", "/dish/find/id", "/dish/find/name", "/dish/find/category").
		IgnorePaths("/category/list", "/category/find/id", "/category/find/name").
		IgnorePaths("/review/find/id", "/review/find/dish_id", "/review/find/rating").
		Build())
	return server
}
