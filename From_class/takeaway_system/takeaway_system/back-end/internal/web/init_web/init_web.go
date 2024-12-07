package init_web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, server *gin.Engine) {
	RegisterCustomerRoutes(db, server)
	RegisterEmployeeRoutes(db, server)
	RegisterDishRoutes(db, server)
	RegisterCategoryRoutes(db, server)
	RegisterCartItemRoutes(db, server)
	RegisterOrderRoutes(db, server)
	RegisterOrderItemRoutes(db, server)
	RegisterOrderStatusHistoryRoutes(db, server)
	RegisterReviewRoutes(db, server)
}

func RegisterCustomerRoutes(db *gorm.DB, server *gin.Engine) {
	c := initCustomer(db)
	cg := server.Group("/customer")
	cg.POST("/signup", c.SignUpCustomer)
	cg.POST("/login", c.LogInCustomer)
	cg.POST("/edit", c.EditCustomer)
	cg.POST("/change_password", c.ChangeCustomerPassword)
	cg.GET("/profile", c.ProfileCustomer)
	cg.GET("/logout", c.LogOutCustomer)
	cga := server.Group("/admin/customer")
	cga.POST("/edit", c.EditCustomerByAdmin)
	cga.GET("/list", c.GetAllCustomers)
	cga.GET("/find/id", c.GetCustomerById)
	cga.GET("/find/name", c.GetCustomerByName) // 模糊查询
	cga.POST("/delete", c.DeleteCustomer)
}

func RegisterEmployeeRoutes(db *gorm.DB, server *gin.Engine) {
	e := initEmployee(db)
	eg := server.Group("/employee")
	eg.POST("/signup", e.SignUpEmployee)
	eg.POST("/login", e.LogInEmployee)
	eg.POST("/edit", e.EditEmployee)
	eg.POST("/change_password", e.ChangeEmployeePassword)
	eg.GET("/profile", e.ProfileEmployee)
	eg.GET("/logout", e.LogOutEmployee)
	ega := server.Group("/admin/employee")
	ega.POST("/edit/all", e.EditEmployeeByAdmin)
	ega.POST("/edit/role", e.EditEmployeeRole)
	ega.POST("/edit/status", e.EditEmployeeStatus)
	ega.GET("/list", e.GetAllEmployees)
	ega.GET("/find/id", e.GetEmployeeById)
	ega.GET("/find/name", e.GetEmployeeByName) // 模糊查询
	ega.GET("/find/role", e.GetEmployeeByRole)
	ega.GET("/find/status", e.GetEmployeeByStatus)
	ega.GET("/find/new", e.GetNewEmployees)
	ega.POST("/delete", e.DeleteEmployee)
}

func RegisterDishRoutes(db *gorm.DB, server *gin.Engine) {
	d := initDish(db)
	dg := server.Group("/dish")
	dg.GET("/list", d.GetAllDishes)
	dg.GET("/find/id", d.GetDishById)
	dg.GET("/find/name", d.GetDishByName) // 模糊查询
	dg.GET("/find/category", d.GetDishByCategory)
	dge := server.Group("/employee/dish")
	dge.POST("/create", d.CreateDish)
	dge.POST("/edit", d.EditDish)
	dge.POST("/delete", d.DeleteDish)
}

func RegisterCategoryRoutes(db *gorm.DB, server *gin.Engine) {
	c := initCategory(db)
	cg := server.Group("/category")
	cg.GET("/list", c.GetAllCategories)
	cg.GET("/find/id", c.GetCategoryById)
	cg.GET("/find/name", c.GetCategoryByName)
	cge := server.Group("/employee/category")
	cge.POST("/create", c.CreateCategory)
	cge.POST("/edit", c.EditCategory)
	cge.POST("/delete", c.DeleteCategory)
}

func RegisterCartItemRoutes(db *gorm.DB, server *gin.Engine) {
	c := initCartItem(db)
	cg := server.Group("/cart")
	cg.POST("/add", c.AddCartItem)                  // 仅顾客
	cg.POST("/edit", c.EditCartItem)                // 仅顾客
	cg.POST("/delete", c.DeleteCartItem)            // 仅顾客
	cg.GET("/list", c.GetCartItemsByCustomerId)     // 仅顾客
	cg.GET("/find/id", c.GetCartItemById)           // 仅顾客
	cg.GET("/clear", c.DeleteCartItemsByCustomerId) // 仅顾客
}

func RegisterOrderRoutes(db *gorm.DB, server *gin.Engine) {
	o := initOrder(db)
	og := server.Group("/order")
	og.POST("/create", o.CreateOrder)                          // 仅顾客
	og.POST("/pay", o.PayTheOrder)                             // 仅顾客
	og.GET("/list", o.GetOrdersByCustomerId)                   // 仅顾客
	og.GET("/find/id", o.GetOrderById)                         // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/find/status", o.GetOrdersByStatus)                // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/find/payment_status", o.GetOrdersByPaymentStatus) // 都可以(对顾客有权限控制, 非本人不可查看)
	og.POST("/cancel", o.CancelTheOrder)                       // 仅顾客
	og.POST("/delete", o.DeleteTheOrder)                       // 仅顾客
	oge := server.Group("/employee/order")
	oge.GET("/find/delivery_person_id", o.GetOrdersByDeliveryPersonId)
	oge.GET("/find/customer_id", o.GetOrdersByCustomerIdByEmployee)
	oge.GET("/list", o.EmployeeGetOrders)
	oge.POST("/cancel", o.CancelTheOrderByEmployee)
	oge.POST("/confirm", o.ConfirmTheOrder)
	oge.POST("/complete", o.MealPreparationCompleted)
	ogd := server.Group("/deliveryman/order")
	ogd.GET("/find/waiting_for_delivery", o.DeliverymanGetOrdersWaitingForDelivery)
	ogd.GET("/find/delivering", o.DeliverymanGetOrdersDelivering)
	ogd.GET("/find/delivered", o.DeliverymanGetOrdersDelivered)
	ogd.POST("/deliver", o.DeliverTheFood)
	ogd.POST("/delivered", o.FoodDelivered)
}

func RegisterOrderItemRoutes(db *gorm.DB, server *gin.Engine) {
	o := initOrderItem(db)
	og := server.Group("/order_item")
	og.GET("/find/id", o.GetOrderItemById)                       // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/find/order_id", o.GetOrderItemsByOrderId)           // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/find/dish_id", o.GetOrderItemsByDishIdByCustomer)   // 仅顾客
	og.GET("/find/review_status", o.GetOrderItemsByReviewStatus) // 仅顾客
	og.GET("/list", o.GetAllOrderItemsByCustomer)                // 仅顾客
	oge := server.Group("/employee/order_item")
	oge.GET("/list", o.GetAllOrderItemsByEmployee)
	oge.GET("/find/dish_id", o.GetOrderItemsByDishId)
}

func RegisterOrderStatusHistoryRoutes(db *gorm.DB, server *gin.Engine) {
	o := initOrderStatusHistory(db)
	og := server.Group("/order_status_history")
	og.GET("/find/id", o.FindOrderStatusHistoryByID)                                 // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/find/order_id", o.FindOrderStatusHistoriesByOrderID)                    // 都可以(对顾客有权限控制, 非本人不可查看)
	og.GET("/list", o.FindOrderStatusHistoriesAllByCustomer)                         // 仅顾客
	og.GET("/find/status", o.FindOrderStatusHistoriesByStatusByCustomer)             // 仅顾客
	og.GET("/find/changed_by_id", o.FindOrderStatusHistoriesByChangedByIDByCustomer) // 仅顾客
	oge := server.Group("/employee/order_status_history")
	oge.GET("/list", o.FindOrderStatusHistoriesAllByEmployee)
	oge.GET("/find/status", o.FindOrderStatusHistoriesByStatusByEmployee)
	oge.GET("/find/changed_by_id", o.FindOrderStatusHistoriesByChangedByIDByEmployee)
}

func RegisterReviewRoutes(db *gorm.DB, server *gin.Engine) {
	r := initReview(db)
	rg := server.Group("/review")
	rg.POST("/create", r.CreateReview)        // 仅顾客
	rg.GET("/list", r.GetReviewsByCustomerId) // 仅顾客
	rg.GET("/find/id", r.GetReviewById)
	rg.GET("/find/dish_id", r.GetReviewsByDishId)
	rg.GET("/find/rating", r.GetReviewsByRating)
	rg.POST("/edit", r.EditReview)     // 仅顾客
	rg.POST("/delete", r.DeleteReview) // 仅顾客
}
