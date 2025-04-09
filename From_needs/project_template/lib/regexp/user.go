package regexp

import regexp "github.com/dlclark/regexp2"

const (
	// 和用'"'比起来，用'`'看起来更清爽 (不用转译)
	idRegexPattern    = `^[a-zA-Z0-9\-]{1,50}$`
	nameRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,12}$`
	emailRegexPattern = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	// emailRegexPattern = `^[0-9a-zA-Z_]{0,19}@[0-9a-zA-Z]{1,13}\.[com,cn,net]{1,3}$` 约束到了.com之类的后缀，此处先暂时不用
	// 密码包含至少一位数字，字母和特殊字符,且长度8-16
	passwordRegexPattern = `^(?![0-9a-zA-Z]+$)[a-zA-Z0-9~!@#$%^&*?_-]{8,16}$`
	phoneRegexPattern    = `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	addressRegexPattern  = `^[\u4e00-\u9fa5a-zA-Z0-9 ]{2,50}$`
)

type UserRegExp struct {
	UserId   *regexp.Regexp
	Name     *regexp.Regexp
	Email    *regexp.Regexp
	Password *regexp.Regexp
	Phone    *regexp.Regexp
	Address  *regexp.Regexp
}

func NewUserRegExp() *UserRegExp {
	idExp := regexp.MustCompile(idRegexPattern, regexp.None)
	nameExp := regexp.MustCompile(nameRegexPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	addressExp := regexp.MustCompile(addressRegexPattern, regexp.None)
	return &UserRegExp{
		UserId:   idExp,
		Name:     nameExp,
		Email:    emailExp,
		Password: passwordExp,
		Phone:    phoneExp,
		Address:  addressExp,
	}
}
