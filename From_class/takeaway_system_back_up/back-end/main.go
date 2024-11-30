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
	db := initDB()
	server := initWebServer()
	init_web.RegisterRoutes(db, server)
	server.Run(":1000")
}

func initDB() *gorm.DB {
	//配置MySQL连接参数
	username := "root"  //账号
	password := "1234"  //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306        //数据库端口
	Dbname := "test"    //数据库名
	timeout := "10s"    //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)

	}

	err = dao.InitTable(db)

	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
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
			return strings.Contains(origin, "szcq.top")
		},
		MaxAge: 12 * time.Second,
	}))
	// 登录校验
	store := cookie.NewStore([]byte("secret"))  // 存 sessions 数据的地方
	server.Use(sessions.Sessions("sid", store)) // cookie 的名字和值(store)
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/customer/login").
		IgnorePaths("/customer/signup").
		IgnorePaths("/employee/login").
		IgnorePaths("/employee/signup").
		IgnorePaths("/dish/*").
		IgnorePaths("/category/*").
		IgnorePaths("/review/find/*").
		Build())
	return server
}
