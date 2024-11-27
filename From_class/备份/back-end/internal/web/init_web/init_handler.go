package init_web

import (
	"back-end/internal/repository"
	"back-end/internal/repository/dao"
	"back-end/internal/servies"
	"back-end/internal/web"
	"gorm.io/gorm"
)

func initCustomer(db *gorm.DB) *web.CustomerHandler {
	cd := dao.NewCustomerDAO(db)
	repo := repository.NewCustomerRepository(cd)
	svc := servies.NewUserService(repo)
	c := web.NewCustomerHandler(svc)
	return c
}
