package lib

import (
	"project/lib/regexp"
)

var UserRegExp *regexp.UserRegExp

func InitRegExp() {
	UserRegExp = regexp.NewUserRegExp()
	if UserRegExp == nil {
		ErrorLog.Fatal("UserRegExp initialization failed")
	}
}
