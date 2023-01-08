package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
)

// UserDataVerify 校验 Type Phone Password 格式
func UserDataVerify(user *model.UserAuthor) *code.MsgCode {
	if user.Type != model.UserTypeCustomer && user.Type != model.UserTypeMerchant {
		return &code.MsgCode{
			Msg:  "Invalid Identity.",
			Code: code.InvalidIdentity,
		}
	} else if !utils.CheckMobile(user.Phone) {
		return &code.MsgCode{Msg: "Invalid Phone.", Code: code.InvalidPhone}
	} else if !utils.CheckPassword(user.Password) {
		return &code.MsgCode{Msg: "Invalid Password.", Code: code.InvalidPassword}
	}
	return &code.MsgCode{
		Msg: "OK", Code: code.OK,
	}
}

func UserCreate(tx *gorm.DB, user *model.UserAuthor) (*model.UserAuthor, *code.MsgCode, error) {
	// 数据校验
	if msgCode := UserDataVerify(user); msgCode.Code != code.OK {
		return nil, msgCode, nil
	}
	// 重复校验
	_, msgCode, err := UserCheck(&model.UserAuthor{
		Phone: user.Phone,
		Type:  user.Type,
	})

	if msgCode.Code == code.CheckError {
		return nil, msgCode, err
	}

	if msgCode.Code == code.OK {
		return nil, &code.MsgCode{Msg: "UserExist", Code: code.UserExist}, nil
	}

	// 生成盐值， 加密密码
	salt := utils.Salt(32)
	pwdMD5, err2 := utils.PasswordMD5(user.Password, salt)

	if err2 != nil {
		return nil, &code.MsgCode{Msg: "ServerError", Code: code.ServerError}, err2
	}

	insertUser := &model.UserAuthor{
		Phone:      user.Phone,
		Type:       user.Type,
		Salt:       salt,
		Password:   pwdMD5,
		PrivateKey: "",
	}

	// 数据库存储
	result := tx.Create(insertUser)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return insertUser, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

//func UserDrop(user *model.UserAuthor) {}

// UserCheck CheckError; DBEmpty; OK;
func UserCheck(user *model.UserAuthor) (*model.UserAuthor, *code.MsgCode, error) {

	searchUser := &model.UserAuthor{}

	// 条件由外部决定
	result := model.Db.Self.Where(user).Find(searchUser)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchUser, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// UserUpdate user 就是要保存的数据
func UserUpdate(tx *gorm.DB, user *model.UserAuthor) (*code.MsgCode, error) {

	result := tx.Save(user)

	if result.Error != nil {
		return &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
