package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RegisterCustomerModel struct {
	Birthday  *string `json:"birthday"`
	Introduce string  `json:"introduce"`
	Nickname  string  `json:"nickname"`
	Password  string  `json:"password"`
	Phone     string  `json:"phone"`
}

// RegisterMerchantModel regist_merchant
type RegisterMerchantModel struct {
	Address       string `json:"address"`
	Introduce     string `json:"introduce"`
	Nickname      string `json:"nickname"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	ShopIntroduce string `json:"shop_introduce"`
	ShopName      string `json:"shop_name"`
}

type LoginModel struct {
	LoginType uint8  `json:"login_type"` // type = 1 客户; type = 2 商家;
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

type MobileUpdateModel struct {
	Password string `json:"password"`
	Phone    string `json:"phone"`
}
type PasswordUpdateModel struct {
	Password    string `json:"password"`
	PasswordNew string `json:"password_new"`
}

func PublicKeyHandler(c *gin.Context) {}

func LoginHandler(c *gin.Context) {
	user := &LoginModel{}
	// 参数绑定失败
	if err := c.BindJSON(user); err != nil {
		code.GinBadRequest(c)
		return
	}

	checkUser, msgCode, _ := dbop.UserCheck(&model.UserAuthor{
		Phone: user.Phone,
		Type:  user.LoginType,
	})

	// 查询失败
	if msgCode.Code != code.OK && msgCode.Code != code.DBEmpty {
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinPhoneORPasswordError(c)
		return
	}

	// 查询成功，开始校验
	pwdMD5, err := utils.PasswordMD5(user.Password, checkUser.Salt)

	if err != nil {
		code.GinServerError(c)
		return
	}

	if pwdMD5 != checkUser.Password {
		code.GinPhoneORPasswordError(c)
		return
	}

	// 验证成功， 开始生成token
	token, err2 := utils.TokenEecode(&utils.TokenPayload{
		Type:  checkUser.Type,
		ID:    checkUser.ID,
		Phone: checkUser.Phone,
	})

	if err2 != nil {
		code.GinServerError(c)
		return
	}

	// 发送token
	code.GinOKPayload(c, &gin.H{
		"token": token,
	})

}
func RegisterCustomerHandler(c *gin.Context) {

	tx := model.Db.Self.Begin()

	customer := &RegisterCustomerModel{}
	// 参数绑定失败
	if err := c.BindJSON(customer); err != nil {
		code.GinBadRequest(c)
		return
	}

	user, msgCode, _ := dbop.UserCreate(tx, &model.UserAuthor{
		Phone:    customer.Phone,
		Password: customer.Password,
		Type:     model.UserTypeCustomer,
	})
	// 插入失败
	if msgCode.Code != code.OK {
		code.GinEmptyMsgCode(c, msgCode)
		return
	}

	// 插入成功 开始插入个人信息
	_, msgCode2, _ := dbop.CustomerInfoCreate(tx, &model.CustomerInfo{
		CustomerID: user.ID,
		NickName:   customer.Nickname,
		Intro:      customer.Introduce,
		Birth:      utils.TimeParse(*customer.Birthday),
	})

	// 插入失败，回滚注册
	if msgCode2.Code != code.OK {
		tx.Rollback()
		code.GinEmptyMsgCode(c, msgCode2)
		return
	}

	// 提交事务
	tx.Commit()
	code.GinOKEmpty(c)
}
func RegisterMerchantHandler(c *gin.Context) {

	tx := model.Db.Self.Begin()

	merchant := &RegisterMerchantModel{}
	// 参数绑定失败
	if err := c.BindJSON(merchant); err != nil {
		code.GinBadRequest(c)
		return
	}

	// 插入之前先校验店铺名
	_, msgCode0, _ := dbop.MerchantInfoCheck(&model.MerchantInfo{
		ShopName: merchant.ShopName,
	})

	// 店铺名重复
	if msgCode0.Code != code.DBEmpty {
		c.JSON(200, gin.H{
			"code": code.ShopNameExist,
			"msg":  "ShopNameExist",
			"data": gin.H{},
		})
		return
	}

	user, msgCode, _ := dbop.UserCreate(tx, &model.UserAuthor{
		Phone:    merchant.Phone,
		Password: merchant.Password,
		Type:     model.UserTypeMerchant,
	})

	// 插入失败
	if msgCode.Code != code.OK {
		code.GinEmptyMsgCode(c, msgCode)
		return
	}

	// 插入成功 开始插入个人信息
	_, msgCode2, _ := dbop.MerchantInfoCreate(tx, &model.MerchantInfo{
		MerchantID: user.ID,
		NickName:   merchant.Nickname,
		Intro:      merchant.Introduce,
		ShopName:   merchant.ShopName,
		ShopIntro:  merchant.ShopIntroduce,
		Address:    merchant.Address,
	})

	// 插入失败，回滚注册
	if msgCode2.Code != code.OK {
		tx.Rollback()
		code.GinEmptyMsgCode(c, msgCode2)
		return
	}

	// 提交事务
	tx.Commit()
	code.GinOKEmpty(c)
}

func ShopnameCheckHandler(c *gin.Context) {
	shopName := c.Query("shop_name")

	_, msgCode, _ := dbop.MerchantInfoCheck(&model.MerchantInfo{
		ShopName: shopName,
	})

	if msgCode.Code == code.DBEmpty {
		code.GinOKPayload(c, &gin.H{
			"exist": false,
		})
		return
	} else if msgCode.Code == code.OK {
		code.GinOKPayload(c, &gin.H{
			"exist": true,
		})
		return
	}

	code.GinEmptyMsgCode(c, msgCode)

}
func MobileCheckHandler(c *gin.Context) {
	phone := c.Query("phone")
	checkType, err := strconv.Atoi(c.Query("type"))

	if err != nil {
		code.GinServerError(c)
		return
	}

	_, msgCode, _ := dbop.UserCheck(&model.UserAuthor{
		Phone: phone,
		Type:  uint8(checkType),
	})

	if msgCode.Code == code.DBEmpty {
		code.GinOKPayload(c, &gin.H{
			"exist": false,
		})
		return
	} else if msgCode.Code == code.OK {
		code.GinOKPayload(c, &gin.H{
			"exist": true,
		})
		return
	}

	code.GinEmptyMsgCode(c, msgCode)

}

func PrivateKeyHandler(c *gin.Context) {}
func MobileUpdateHandler(c *gin.Context) {

	mobileUpdateModel := &MobileUpdateModel{}

	if err := c.BindJSON(mobileUpdateModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	ID, exist := c.Get("ID")
	//    Type, exist2 := c.Get("ID")
	if !exist {
		code.GinUnAuthorized(c)
	}
	uintID, _ := ID.(uint)
	//    uintType, _ := Type.(uint8)

	user := &model.UserAuthor{
		ID: uintID,
	}

	// 查找 token 的 所有者
	checkUser, msgCode, _ := dbop.UserCheck(user)

	if msgCode.Code == code.CheckError {
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUnAuthorized(c)
		return
	}

	// 校验密码
	pwdMd5, err := utils.PasswordMD5(mobileUpdateModel.Password, checkUser.Salt)
	if err != nil {
		code.GinServerError(c)
		return
	}

	if pwdMd5 != checkUser.Password {
		code.GinPhoneORPasswordError(c)
		return
	}

	// 检查手机号是否重复
	_, msgCode2, _ := dbop.UserCheck(&model.UserAuthor{
		Type:  checkUser.Type,
		Phone: mobileUpdateModel.Phone,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.OK {
		c.JSON(200, gin.H{
			"code": code.PhoneExist,
			"msg":  "PhoneExist",
			"data": gin.H{},
		})
		return
	}
	// 手机号不重复
	checkUser.Phone = mobileUpdateModel.Phone
	if codeMsg3, _ := dbop.UserUpdate(model.Db.Self, checkUser); codeMsg3.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}
	code.GinOKEmpty(c)

}
func PasswordUpdateHandler(c *gin.Context) {

	passwordUpdateModel := &PasswordUpdateModel{}

	if err := c.BindJSON(passwordUpdateModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	ID, exist := c.Get("ID")
	if !exist {
		code.GinUnAuthorized(c)
	}
	uintID, _ := ID.(uint)

	// 查找 token 的 所有者
	checkUser, msgCode, _ := dbop.UserCheck(&model.UserAuthor{
		ID: uintID,
	})

	if msgCode.Code == code.CheckError {
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUnAuthorized(c)
		return
	}

	// 校验密码
	pwdMd5, err := utils.PasswordMD5(passwordUpdateModel.Password, checkUser.Salt)
	if err != nil {
		code.GinServerError(c)
		return
	}

	if pwdMd5 != checkUser.Password {
		code.GinPhoneORPasswordError(c)
		return
	}

	// 修改密码
	newPwdMd5, err2 := utils.PasswordMD5(passwordUpdateModel.PasswordNew, checkUser.Salt)
	if err2 != nil {
		code.GinServerError(c)
		return
	}

	checkUser.Password = newPwdMd5

	if codeMsg3, _ := dbop.UserUpdate(model.Db.Self, checkUser); codeMsg3.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}
	code.GinOKEmpty(c)
}
