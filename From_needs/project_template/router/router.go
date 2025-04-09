package router

import (
	"github.com/gin-gonic/gin"
	"project/controller"
	"project/lib"
)

func SetupRoute(cfg *lib.Config) *gin.Engine {
	router := InitServer(cfg)
	userRoute(router.Group("/user"))
	return router
}

func userRoute(user *gin.RouterGroup) {
	user.POST("/signup", controller.SignUp)
	user.POST("/login", controller.Login)
	user.GET("/logout", controller.LogOut)
	user.GET("/profile", controller.Profile)
	user.POST("/edit/profile", controller.ChangeProfile)
	user.POST("/edit/password", controller.ChangePassword)
}
