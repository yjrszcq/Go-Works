package init_web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, server *gin.Engine) {
	RegisterCustomerRoutes(db, server)
}

func RegisterCustomerRoutes(db *gorm.DB, server *gin.Engine) {
	c := initCustomer(db)
	// 分组 (以防"/users"这个前缀手一抖写错)
	cg := server.Group("/customer")
	cg.POST("/signup", c.SignUp)
	cg.POST("/login", c.LogIn)
	cg.POST("/edit", c.Edit)
	cg.GET("/profile", c.Profile)
}
