package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
)

type CommodityCreateModel struct {
	Count   uint    `json:"count"` // 商品库存
	Intro   string  `json:"intro"`
	Name    string  `json:"name"`
	Picture string  `json:"picture"`
	Price   float64 `json:"price"` // 商品单价
}

type CommodityUpdateModel struct {
	Count   uint    `json:"count"` // 商品库存
	ID      uint    `json:"id"`
	Intro   string  `json:"intro"`
	Name    string  `json:"name"`
	Picture string  `json:"picture"`
	Price   float64 `json:"price"` // 商品单价
}

// CommodityStatusDecide 商品数量不足， 但是也是上架状态
func CommodityStatusDecide(current uint, status bool, count uint) uint {

	var minCount uint = utils.GlobalConfig.Global.GroupMemberCount

	// 强制下架 和 删除 无法更改
	if current == model.CommodityStatusForceOnShelf || current == model.CommodityStatusDeleted {
		return current
	}

	if status {
        // 想上架， 看数量
        if count >= minCount {
            return model.CommodityStatusOnShelf
        } else if count > 0 {
            return model.CommodityStatusNotEnought
        }
	}
	return model.CommodityStatusOffShelf
}

func CommodityCreateHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	commodityCreateModel := &CommodityCreateModel{}
	if !utils.QuickBind(c, commodityCreateModel) {
		return
	}

	if commodityCreateModel.Count <= 0 || commodityCreateModel.Price <= 0 {
		code.GinBadRequest(c)
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCreate(model.Db.Self, &model.CommodityInfo{
		MerchantID: ID,
		Count:      commodityCreateModel.Count,
		Name:       commodityCreateModel.Name,
		Price:      commodityCreateModel.Price,
		Intro:      commodityCreateModel.Intro,
		Status:     CommodityStatusDecide(model.CommodityStatusNotCreate, false, commodityCreateModel.Count),
		Picture:    commodityCreateModel.Picture,
	})

	if msgCode.Code == code.InsertError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"id":      commodity.ID,
		"status":  commodity.Status,
		"name":    commodity.Name,
		"intro":   commodity.Intro,
		"picture": commodity.Picture,
		"price":   commodity.Price,
		"count":   commodity.Count,
	})
}
func CommodityUpdateHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	commodityUpdateModel := &CommodityUpdateModel{}
	if !utils.QuickBind(c, commodityUpdateModel) {
		return
	}

	if commodityUpdateModel.Count <= 0 || commodityUpdateModel.Price <= 0 {
		code.GinBadRequest(c)
		return
	}

	_, msgCode, _ := dbop.CommodityInfoUpdate(
		model.Db.Self,
		&model.CommodityInfo{
			MerchantID: ID,
			ID:         commodityUpdateModel.ID,
		},
		&model.CommodityInfo{
			Count:   commodityUpdateModel.Count,
			Name:    commodityUpdateModel.Name,
			Price:   commodityUpdateModel.Price,
			Intro:   commodityUpdateModel.Intro,
			Picture: commodityUpdateModel.Picture,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	code.GinOKEmpty(c)
}

type CommodityStatusModel struct {
	ID     uint `json:"commodity_id"` // 商品id
	Status bool `json:"status"`       // true 上架
}

func CommodityStatusHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	commodityStatusModel := &CommodityStatusModel{}
	if !utils.QuickBind(c, commodityStatusModel) {
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		MerchantID: ID,
		ID:         commodityStatusModel.ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUnMatchedID(c)
		return
	}

	if commodity[0].Status != model.CommodityStatusOnShelf && commodity[0].Status != model.CommodityStatusOffShelf {
		code.GinOKPayload(c, &gin.H{
			"status": commodity[0].Status,
		})
		return
	}

	finalStatus := CommodityStatusDecide(commodity[0].Status, commodityStatusModel.Status, commodity[0].Count)

	updateCommodity, msgCode2, _ := dbop.CommodityInfoUpdate(
		model.Db.Self,
		&model.CommodityInfo{
			MerchantID: ID,
			ID:         commodityStatusModel.ID,
		},
		&model.CommodityInfo{
			Status: finalStatus,
		})

	if msgCode2.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinOKPayload(c, &gin.H{
			"status": finalStatus,
		})
		return
	}

	code.GinOKPayload(c, &gin.H{
		"status": updateCommodity.Status,
	})
}

func CommodityDeleteHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		ID:         pathIDModel.ID,
		MerchantID: ID,
	})
	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	if commodity[0].Status == model.CommodityStatusOnShelf {
		code.GinOKPayload(c, &gin.H{
			"status": commodity[0].Status,
		})
		return
	}

	_, msgCode2, _ := dbop.CommodityInfoUpdate(
		model.Db.Self,
		&model.CommodityInfo{
			ID:         pathIDModel.ID,
			MerchantID: ID,
		},
		&model.CommodityInfo{
			Status: model.CommodityStatusDeleted,
		},
	)

	if msgCode2.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}
	code.GinOKPayload(c, &gin.H{
		"status": model.CommodityStatusDeleted,
	})
	return
}
func CommodityDetailHandler(c *gin.Context) {
	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		ID: pathIDModel.ID,
	})
	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	code.GinOKPayloadAny(c, commodity[0])
}

func CommodityShopListHandler(c *gin.Context) {

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	limit := c.Query("limit")
	page := c.Query("page")

	commoditys, msgCode, _ := dbop.CommodityInfoLimitPageCheck(&model.CommodityInfo{
		MerchantID: pathIDModel.ID,
	}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"list":  commoditys,
		"count": len(commoditys),
	})

}
