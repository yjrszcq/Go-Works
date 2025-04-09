package service

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"project/domain"
	"project/lib"
	"project/model"
)

var (
	ErrorUserIsUnavailable          = errors.New("用户不可用")
	ErrorUserHasNoPermission        = errors.New("无权限")
	ErrorUserDuplicateEmail         = model.ErrorUserDuplicateEmail
	ErrorUserInvalidUserOrPassword  = errors.New("邮箱或密码错误")
	ErrorUserNotFound               = model.ErrorUserNotFound
	ErrorUserPasswordIsWrong        = errors.New("密码错误")
	ErrorUserPasswordIsInconsistent = errors.New("两次输入的密码不一致")
	ErrorUserFormatOfName           = errors.New("用户名格式错误")
	ErrorUserFormatOfEmail          = errors.New("邮箱格式错误")
	ErrorUserFormatOfPassword       = errors.New("密码格式错误")
)

func getCurrentSession(ctx *gin.Context) (sessions.Session, error) {
	sess := sessions.Default(ctx)
	if sess == nil || sess.Get("status").(string) == "unavailable" {
		return nil, ErrorUserIsUnavailable
	}
	return sess, nil
}

func SignUp(ctx *gin.Context, u *domain.User) (*domain.User, error) {
	ok, err := lib.UserRegExp.Email.MatchString(u.Email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrorUserFormatOfEmail
	}
	if u.ConfirmPassword != u.Password {
		return nil, ErrorUserPasswordIsInconsistent
	}
	ok, err = lib.UserRegExp.Password.MatchString(u.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrorUserFormatOfPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)
	id, err := model.InsertUser(ctx, &model.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		if errors.Is(err, model.ErrorUserDuplicateEmail) {
			return nil, ErrorUserDuplicateEmail
		} else {
			return nil, err
		}
	}
	return &domain.User{
		UserId: id,
	}, nil
}

func LogIn(ctx *gin.Context, u *domain.User) (*domain.User, error) {
	// 先找用户
	user, err := model.SelectUser(ctx, model.User{Email: u.Email})
	if err != nil {
		if errors.Is(err, model.ErrorUserNotFound) {
			return nil, ErrorUserInvalidUserOrPassword
		} else {
			return nil, err
		}
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return nil, ErrorUserInvalidUserOrPassword
	}
	// 设置 session
	sess := sessions.Default(ctx)
	sess.Set("id", user.UserId)
	sess.Set("status", "available")
	err = sess.Save()
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &domain.User{
		UserId: user.UserId,
	}, nil
}

func LogOut(ctx *gin.Context) error {
	sess := sessions.Default(ctx)
	sess.Clear()
	err := sess.Save()
	if err != nil {
		return err
	}
	return nil
}

func Profile(ctx *gin.Context) (*domain.User, error) {
	session, err := getCurrentSession(ctx)
	if err != nil {
		return nil, err
	}
	user, err := model.SelectUser(ctx, model.User{UserId: session.Get("id").(string)})
	if err != nil {
		if errors.Is(err, model.ErrorUserNotFound) {
			return nil, ErrorUserInvalidUserOrPassword
		} else {
			return nil, err
		}
	}
	return &domain.User{
		UserId: user.UserId,
		Name:   user.Name,
		Email:  user.Email,
	}, err
}

func ChangeProfile(ctx *gin.Context, u *domain.User) error {
	session, err := getCurrentSession(ctx)
	if err != nil {
		return err
	}
	ok, err := lib.UserRegExp.Name.MatchString(u.Name)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorUserFormatOfName
	}
	ok, err = lib.UserRegExp.Name.MatchString(u.Email)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorUserFormatOfEmail
	}
	err = model.UpdateUser(ctx, model.User{
		UserId: session.Get("id").(string),
		Name:   u.Name,
		Email:  u.Email,
	})
	if err != nil {
		if errors.Is(err, model.ErrorUserDuplicateEmail) {
			return ErrorUserDuplicateEmail
		} else if errors.Is(err, model.ErrorUserNotFound) {
			return ErrorUserInvalidUserOrPassword
		} else {
			return err
		}
	}
	return nil
}

func ChangePassword(ctx *gin.Context, u *domain.User) error {
	session, err := getCurrentSession(ctx)
	if err != nil {
		return err
	}
	if u.ConfirmPassword != u.NewPassword {
		return ErrorUserPasswordIsInconsistent
	}
	ok, err := lib.UserRegExp.Password.MatchString(u.NewPassword)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorUserFormatOfPassword
	}
	user, err := model.SelectUser(ctx, model.User{UserId: session.Get("id").(string)})
	if err != nil {
		if errors.Is(err, model.ErrorUserNotFound) {
			return ErrorUserNotFound
		} else {
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return ErrorUserPasswordIsWrong
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newPassword := string(hash)
	err = model.UpdateUserPassword(ctx, model.User{
		UserId:   user.UserId,
		Password: newPassword,
	})
	if err != nil {
		return err
	}
	return nil
}
