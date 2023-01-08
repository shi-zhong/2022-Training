package utils

import (
	"regexp"
	"time"
)

func CheckMobile(mobile string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(mobile)
}

func HideMobile(mobile string) string {
	return mobile[:3] + "****" + mobile[7:]
}

func CheckPassword(pwd string) bool {

	regRuler := "^[0-9A-Za-z!@#$%^&*]{8,16}$"
	reg := regexp.MustCompile(regRuler)

	return reg.MatchString(pwd)
}

func TimeParse(t string) time.Time {
	tt, _ := time.Parse("2006-01-02", t)
	return tt
}
