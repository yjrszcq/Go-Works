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
	return c
}

func initEmployee(db *gorm.DB) *web.EmployeeHandler {
	ed := dao.NewEmployeeDAO(db)
	repo := repository.NewEmployeeRepository(ed)
	svc := service.NewEmployeeService(repo)
	e := web.NewEmployeeHandler(svc)
	return e
}
