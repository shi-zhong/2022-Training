package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
)

func CartListHandler(c *gin.Context) {

	ID, _, _ := utils.GetTokenInfo(c)

	limit := c.Query("limit")
	page := c.Query("page")

	carts, msgCode, _ := dbop.ShoppingCartLimitPageUnionCheck(&model.ShoppingCart{
		CustomerID: ID,
	}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"list":  carts,
		"count": len(carts),
	})

}
func CartListAddHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	commodity, msgCode2, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		ID: pathIDModel.ID,
	})
	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	if commodity[0].Status != model.CommodityStatusOnShelf {
		code.GinBadRequest(c)
		return
	}

	_, msgCode, _ := dbop.ShoppingCartCreate(model.Db.Self, &model.ShoppingCart{
		CustomerID:  ID,
		CommodityID: pathIDModel.ID,
	})

	if msgCode.Code == code.InsertError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
func CartListRemoveHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	msgCode, _ := dbop.ShoppingCartDrop(model.Db.Self, &model.ShoppingCart{
		CustomerID: ID,
		ID:         pathIDModel.ID,
	})

	if msgCode.Code == code.DropError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
