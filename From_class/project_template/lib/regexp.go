package lib

import (
	"GoLandProjects/Works/From_class/haze_detection_system/lib/regexp"
)

var UserRegExp *regexp.UserRegExp

func InitRegExp() {
	UserRegExp = regexp.NewUserRegExp()
	if UserRegExp == nil {
		ErrorLog.Fatal("UserRegExp initialization failed")
	}
}
