package dbop
//
//import (
//	"HappyShopTogether/model"
//	"HappyShopTogether/utils"
//	"HappyShopTogether/utils/code"
//)
//
//func AdminDataVerify(admin *model.AdministratorAuthor) (string, int) {
//	if admin.AdminAuthor.Type != 3 && admin.AdminAuthor.Type != 2 {
//		return "Invalid Idnetity.", code.InvalidIdentity
//	} else if !utils.CheckMobile(admin.AdminAuthor.Phone) {
//		return "Invalid Phone.", code.InvalidPhone
//	} else if !utils.CheckPassword(admin.AdminAuthor.Password) {
//		return "Invalid Password.", code.InvalidPassword
//	}
//	return "", 0
//}
//
//func AdminCreate(admin *model.AdministratorAuthor) (string, int, error) {
//	return "OK", code.OK, nil
//}
//
//func AdminDrop(admin *model.AdministratorAuthor) (string, int, error) {
//	return "", 0, nil
//}
//
//func AdminCheck(admin *model.AdministratorAuthor) (string, int, error) {
//	return "", 0, nil
//}
//
//func AdminUpdate(admin *model.AdministratorAuthor) (string, int, error) {
//	return "", 0, nil
//}
