package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UpdateCustomerModel struct {
	Avatar    *string `json:"avatar"`
	Introduce *string `json:"introduce"`
	Nickname  *string `json:"nickname"`
}

type UpdateMerchantModel struct {
	Address    string  `json:"address"`
	Avatar     *string `json:"avatar"`
	Introduce  *string `json:"introduce"`
	Nickname   string  `json:"nickname"`
	ShopAvatar *string `json:"shop_avatar"`
	ShopIntro  *string `json:"shop_intro"`
}

func GetCustomerHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	queryID := c.Query("user_id")

	var searchID uint
	if queryID == "" {
		searchID = ID
	} else {
		queryID2, err := strconv.Atoi(queryID)
		if err != nil {
			code.GinServerError(c)
			return
		}
		searchID = uint(queryID2)
	}

	// search
	customerInfo, msgCode, _ := dbop.CustomerInfoCheck(&model.CustomerInfo{
		CustomerID: searchID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	// search Phone
	customer, msgCode2, _ := dbop.UserCheck(&model.UserAuthor{
		ID: searchID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"nickname":  customerInfo.NickName,
		"introduce": customerInfo.Intro,
		"phone":     utils.HideMobile(customer.Phone),
		"avatar":    customerInfo.Avatar,
		"birthday":  customerInfo.Birth,
		"user_id":   customerInfo.CustomerID,
	})
}
func GetMerchantHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	queryID := c.Query("merchant_id")

	var searchID uint
	if queryID == "" {
		searchID = ID
	} else {
		queryID2, err := strconv.Atoi(queryID)
		if err != nil {
			code.GinServerError(c)
			return
		}
		searchID = uint(queryID2)
	}

	// search
	merchantInfo, msgCode, _ := dbop.MerchantInfoCheck(&model.MerchantInfo{
		MerchantID: searchID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	// search Phone
	merchant, msgCode2, _ := dbop.UserCheck(&model.UserAuthor{
		ID: searchID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"nickname":    merchantInfo.NickName,
		"introduce":   merchantInfo.Intro,
		"phone":       utils.HideMobile(merchant.Phone),
		"avatar":      merchantInfo.Avatar,
		"shop_name":   merchantInfo.ShopName,
		"shop_intro":  merchantInfo.ShopIntro,
		"id":          merchantInfo.MerchantID,
		"shop_avatar": merchantInfo.ShopAvatar,
		"address":     merchantInfo.Address,
	})
}
func UpdateCustomerHandler(c *gin.Context) {

	ID, _, _ := utils.GetTokenInfo(c)

	updateCustomerModel := &UpdateCustomerModel{}
	if err := c.BindJSON(updateCustomerModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	_, msgCode, _ := dbop.CustomerInfoUpdate(
		model.Db.Self,
		&model.CustomerInfo{
			CustomerID: ID,
		},
		&model.CustomerInfo{
			NickName: *updateCustomerModel.Nickname,
			Intro:    *updateCustomerModel.Introduce,
			Avatar:   *updateCustomerModel.Avatar,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}
	code.GinOKEmpty(c)

	return
}
func UpdateMerhantHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	updateMerchantModel := &UpdateMerchantModel{}
	if err := c.BindJSON(updateMerchantModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	_, msgCode, _ := dbop.MerchantInfoUpdate(
		model.Db.Self,
		&model.MerchantInfo{
			MerchantID: ID,
		},
		&model.MerchantInfo{
			Address:    updateMerchantModel.Address,
			Avatar:     *updateMerchantModel.Avatar,
			Intro:      *updateMerchantModel.Introduce,
			NickName:   updateMerchantModel.Nickname,
			ShopAvatar: *updateMerchantModel.ShopAvatar,
			ShopIntro:  *updateMerchantModel.ShopIntro,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}
	code.GinOKEmpty(c)

	return

}
