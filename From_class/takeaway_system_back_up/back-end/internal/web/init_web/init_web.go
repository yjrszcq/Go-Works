package init_web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, server *gin.Engine) {
	RegisterCustomerRoutes(db, server)
	RegisterEmployeeRoutes(db, server)
}

func RegisterCustomerRoutes(db *gorm.DB, server *gin.Engine) {
	c := initCustomer(db)
	// 分组 (以防"/users"这个前缀手一抖写错)
	cg := server.Group("/customer")
	cg.POST("/signup", c.SignUpCustomer)
	cg.POST("/login", c.LogInCustomer)
	cg.POST("/edit", c.EditCustomer)
	cg.GET("/profile", c.ProfileCustomer)
	cg.GET("/logout", c.LogOutCustomer)
	cg.GET("/list", c.GetAllCustomers)
	cg.GET("/find", c.GetCustomerById)
	cg.POST("/delete", c.DeleteCustomer)
}

func RegisterEmployeeRoutes(db *gorm.DB, server *gin.Engine) {
	e := initEmployee(db)
	eg := server.Group("/employee")
	eg.POST("/signup", e.SignUpEmployee)
	eg.POST("/login", e.LogInEmployee)
	eg.POST("/edit", e.EditEmployee)
	eg.GET("/profile", e.ProfileEmployee)
	eg.GET("/logout", e.LogOutEmployee)
	eg.GET("/list", e.GetAllEmployees)
	eg.GET("/find", e.GetEmployeeById)
	eg.POST("/delete", e.DeleteEmployee)
}
