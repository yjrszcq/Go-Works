package init_web

import (
	"back-end/internal/repository"
	"back-end/internal/repository/dao"
	"back-end/internal/service"
	"back-end/internal/web"
	"gorm.io/gorm"
)

func initCustomer(db *gorm.DB) *web.CustomerHandler {
	cd := dao.NewCustomerDAO(db)
	repo := repository.NewCustomerRepository(cd)
	svc := service.NewCustomerService(repo)
	c := web.NewCustomerHandler(svc)
	service.GlobalCustomer = repo
	return c
}

func initEmployee(db *gorm.DB) *web.EmployeeHandler {
	ed := dao.NewEmployeeDAO(db)
	repo := repository.NewEmployeeRepository(ed)
	svc := service.NewEmployeeService(repo)
	e := web.NewEmployeeHandler(svc)
	service.GlobalEmployee = repo
	return e
}

func initDish(db *gorm.DB) *web.DishHandler {
	dd := dao.NewDishDAO(db)
	repo := repository.NewDishRepository(dd)
	svc := service.NewDishService(repo)
	d := web.NewDishHandler(svc)
	service.GlobalDish = repo
	return d
}

func initCategory(db *gorm.DB) *web.CategoryHandler {
	cd := dao.NewCategoryDAO(db)
	repo := repository.NewCategoryRepository(cd)
	svc := service.NewCategoryService(repo)
	c := web.NewCategoryHandler(svc)
	service.GlobalCategory = repo
	return c
}

func initCartItem(db *gorm.DB) *web.CartItemHandler {
	cid := dao.NewCartItemDAO(db)
	repo := repository.NewCartItemRepository(cid)
	svc := service.NewCartItemService(repo)
	c := web.NewCartItemHandler(svc)
	service.GlobalCartItem = repo
	return c
}

func initOrder(db *gorm.DB) *web.OrderHandler {
	od := dao.NewOrderDAO(db)
	repo := repository.NewOrderRepository(od)
	svc := service.NewOrderService(repo)
	o := web.NewOrderHandler(svc)
	service.GlobalOrder = repo
	return o
}

func initOrderItem(db *gorm.DB) *web.OrderItemHandler {
	oid := dao.NewOrderItemDAO(db)
	repo := repository.NewOrderItemRepository(oid)
	svc := service.NewOrderItemService(repo)
	o := web.NewOrderItemHandler(svc)
	service.GlobalOrderItem = repo
	return o
}

func initOrderStatusHistory(db *gorm.DB) *web.OrderStatusHistoryHandler {
	osd := dao.NewOrderStatusHistoryDAO(db)
	repo := repository.NewOrderStatusHistoryRepository(osd)
	svc := service.NewOrderStatusHistoryService(repo)
	o := web.NewOrderStatusHistoryHandler(svc)
	service.GlobalOrderStatusHistory = repo
	return o
}

func initReview(db *gorm.DB) *web.ReviewHandler {
	rd := dao.NewReviewDAO(db)
	repo := repository.NewReviewRepository(rd)
	svc := service.NewReviewService(repo)
	r := web.NewReviewHandler(svc)
	service.GlobalReview = repo
	return r
}
