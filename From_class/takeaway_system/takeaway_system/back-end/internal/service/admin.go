package service

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminService struct {
	UserName string
	Password string
}

var (
	ErrPasswordIsWrongInAdmin = errors.New("用户名或密码错误")
)

func NewAdminService(userName string, password string) *AdminService {
	return &AdminService{
		UserName: userName,
		Password: password,
	}
}

func (svc *AdminService) LogInAdmin(ctx *gin.Context, userName string, password string) error {
	if userName != GlobalAdmin.UserName || password != GlobalAdmin.Password {
		return ErrPasswordIsWrongInAdmin
	}
	sess := sessions.Default(ctx)
	sess.Set("role", "admin")
	sess.Set("id", 0)
	sess.Set("status", "available")
	err := sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func (svc *AdminService) LogOutAdmin(ctx *gin.Context) error {
	sess := sessions.Default(ctx)
	sess.Clear()
	err := sess.Save()
	if err != nil {
		return err
	}
	return nil
}
