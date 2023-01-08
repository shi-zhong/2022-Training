package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type ShareBillCreateModel struct {
	ModeID    uint `json:"mode_id"`
	AddressID uint `json:"address_id"`
	Mode      bool `josn:"mode"`
}

type ShareBillJoinModel struct {
	ShareBillID string `json:"share_bill_id"`
	OwnerID     uint   `json:"owner_id"`
	AddressID   uint   `json:"address_id"`
}

// ShareBillDetailsModel share_bill
type ShareBillDetailsModel struct {
	CommodityInfo *model.CommodityInfo `json:"commodity_info"`
	CreateAt      time.Time            `json:"create_at"`
	DoneAt        time.Time            `json:"done_at"`
	Member        []*dbop.Member       `json:"member"`
	OwnerID       uint                 `json:"owner_id"`
	ShareBillID   string               `json:"share_bill_id"`
	ShopName      string               `json:"shop_name"`
	ShopAvatar    string               `json:"shop_avatar"`
	Status        uint8                `json:"status"`
}

func beforeCreateCheckCommodity(c *gin.Context, condition *model.CommodityInfo) (*model.CommodityInfo, bool) {
	commodity, msgCode, _ := dbop.CommodityInfoCheck(condition)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingItems(c)
		return nil, false
	}

	return commodity[0], true
}

func beforeCreateCheckAddress(c *gin.Context, condition *model.CustomerAddress) (*model.CustomerAddress, bool) {
	address, msgCode, _ := dbop.CustomerAddressCheck(condition)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode.Code == code.DBEmpty {
		code.GinUnMatchedID(c)
		return nil, false
	}

	if address[0].Default == model.AddressDelete {
		code.GinMissingAddress(c)
		return nil, false
	}

	return nil, true
}

func afterCreateCommodityUpdate(c *gin.Context, tx *gorm.DB, condition, update *model.CommodityInfo) bool {
	_, msgCode, _ := dbop.CommodityInfoUpdate(tx, condition, update)

	if msgCode.Code == code.UpdateError {
		tx.Rollback()
		code.GinServerError(c)
		return false
	} else if msgCode.Code == code.DBEmpty {
		tx.Rollback()
		code.GinUnMatchedID(c)
		return false
	}
	return true
}

/*
   6. 计时器 规定时间后执行
   (1) 查询拼单信息 若成团, 结束
   (2) 未成团, 更新拼单信息
   (3) 查找成员 更新对应订单信息
   (4) 返还商品数量
   (5) 若失败，一段时间后重新执行
*/

func shareBillTimeOverCheck(ID string) bool {
	afterShareBill, msgCode, _ := dbop.ShareBillCheck(&model.ShareBill{ID: ID})
	// 查不到 或者 成团 就结束
	if msgCode.Code == code.DBEmpty || afterShareBill[0].Status == model.ShareBillSuccess {
		return true
	}

	tx := model.Db.Self.Begin()

	// (2) 未成团, 更新拼单信息
	_, _, err := dbop.ShareBillUpdate(tx, &model.ShareBill{
		ID: ID,
	}, &model.ShareBill{
		Status:   model.ShareBillFailed,
		FinishAt: time.Now(),
	})

	if err != nil {
		return false
	}

	// (3) 查找成员 更新对应订单信息

	// get
	orders, msgCode2, _ := dbop.OrderCheck(&model.Order{
		ShareBillID: ID,
	})

	if msgCode2.Code == code.CheckError || msgCode2.Code == code.DBEmpty {
		tx.Rollback()
		return false
	}
	// update  拼团失败
	for _, value := range orders {
		_, msgCode9, _ := dbop.OrderUpdate(tx, &model.Order{ID: value.ID}, &model.Order{
			Status:   model.OrderCancel,
			FinishAt: time.Now()})

		if msgCode9.Code == code.UpdateError {
			tx.Rollback()
			return false
		}
	}
	// (4) 返还商品数量
	commodity, msgCode3, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{ID: afterShareBill[0].CommodityID})

	if msgCode3.Code == code.CheckError {
		tx.Rollback()
		return false
	} else if msgCode3.Code == code.DBEmpty {
		tx.Rollback()
		return false
	}

	_, msgCode4, _ := dbop.CommodityInfoUpdate(tx,
		&model.CommodityInfo{ID: afterShareBill[0].CommodityID},
		&model.CommodityInfo{
			Count:  commodity[0].Count + uint(len(orders)),
			Status: CommodityStatusDecide(commodity[0].Status, true, commodity[0].Count+uint(len(orders))),
		})

	if msgCode4.Code == code.UpdateError {
		tx.Rollback()
		return false
	} else if msgCode4.Code == code.DBEmpty {
		tx.Rollback()
		return false
	}

	tx.Commit()

	return true
}

// 执行拼团结束代码， 若失败则5s后重新执行
func timerOverRepeatCheck(ID string) {
	if !shareBillTimeOverCheck(ID) {
		time.AfterFunc(5*time.Second, func() {
			timerOverRepeatCheck(ID)
		})
	}
}

// ShareBillCreateHandler 创建拼单订单
/*
   before
   1. 传入 地址id 商品id 个人id， 分别校验

   2. 创建 拼单信息
   3. 主人入团
   4. 创建个人订单

   after
   5. 更新商品数量
   6. 计时器 规定时间后执行
       (1) 查询拼单信息 若成团, 结束
       (2) 未成团, 更新拼单信息
       (3) 查找成员 更新对应订单信息
       (4) 返还商品数量
       (5) 若失败，一段时间后重新执行

*/
func ShareBillCreateHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	shareBillCreateModel := &ShareBillCreateModel{}
	if !utils.QuickBind(c, shareBillCreateModel) {
		return
	}

	var cart *model.ShoppingCart

	// in cart
	if shareBillCreateModel.Mode {
		cart0, msgCode0, _ := dbop.ShoppingCartCheck(&model.ShoppingCart{ID: shareBillCreateModel.ModeID, CustomerID: ID})
		if msgCode0.Code == code.CheckError {
			code.GinServerError(c)
			return
		} else if msgCode0.Code == code.DBEmpty {
			code.GinMissingCart(c)
			return
		}
		cart = cart0[0]
	}

	var commodityID uint

	if shareBillCreateModel.Mode {
		commodityID = cart.CommodityID
	} else {
		commodityID = shareBillCreateModel.ModeID
	}

	// 1. 传入 地址id 商品id 个人id， 分别校验
	// 校验商品数量
	commodity, err := beforeCreateCheckCommodity(c, &model.CommodityInfo{ID: commodityID})
	if !err {
		return
	}

	if commodity.Status != model.CommodityStatusOnShelf && commodity.Status != model.CommodityStatusNotEnought {
		code.GinNotOnShelf(c)
		return
	} else if commodity.Status == model.CommodityStatusNotEnought {
		code.GinNotEnough(c)
		return
	}

	// 校验地址
	_, err2 := beforeCreateCheckAddress(c, &model.CustomerAddress{
		ID:         shareBillCreateModel.AddressID,
		CustomerID: ID,
	})

	if !err2 {
		return
	}

	// 校验完毕

	tx := model.Db.Self.Begin()

	// 2. 创建 拼单信息

	shareBill, msgCode2, _ := dbop.ShareBillCreate(tx, &model.ShareBill{
		ID:          utils.ShareBillIDGenerate(),
		OwnerID:     ID,
		CommodityID: commodity.ID,
		Status:      model.ShareBillWaitingForTwo,
		FinishAt:    time.Now(),
		Price:       commodity.Price,
	})

	if msgCode2.Code == code.InsertError || msgCode2.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 3. 主人入团

	_, msgCode3, _ := dbop.ShareBillTeamCreate(tx, &model.ShareBillTeam{
		ShareBillID: shareBill.ID,
		MemberID:    ID,
	})

	if msgCode3.Code == code.InsertError || msgCode3.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 4. 创建个人订单
	order, msgCode5, _ := dbop.OrderCreate(tx, &model.Order{
		ID:          utils.OrderIDGenerate(),
		ShareBillID: shareBill.ID,
		CustomerID:  ID,
		AddressID:   shareBillCreateModel.AddressID,
		DueAt:       time.Now(),
		CommodityAt: time.Now(),
		FinishAt:    time.Now(),
		Status:      model.OrderCreated,
	})

	if msgCode5.Code == code.InsertError || msgCode5.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 5. 更新商品数量
	if !afterCreateCommodityUpdate(c, tx, &model.CommodityInfo{ID: commodity.ID}, &model.CommodityInfo{
		Count:  commodity.Count - 1,
		Status: CommodityStatusDecide(commodity.Status, true, commodity.Count-1),
	}) {
		return
	}

	// 删除购物车
	if shareBillCreateModel.Mode {
		msgCode6, _ := dbop.ShoppingCartDrop(tx, &model.ShoppingCart{ID: shareBillCreateModel.ModeID, CustomerID: ID})
		if msgCode6.Code == code.DropError {
			tx.Rollback()
			code.GinServerError(c)
			return
		} else if msgCode6.Code == code.DBEmpty {
			tx.Rollback()
			code.GinMissingCart(c)
			return
		}
	}

	tx.Commit()

	time.AfterFunc(time.Duration(utils.GlobalConfig.Global.GroupForMemberTime)*time.Hour, func() {
		// 规定时间后执行
		timerOverRepeatCheck(shareBill.ID)
	})

	code.GinOKPayload(c, &gin.H{
		"share_bill_id": shareBill.ID,
		"order_id":      order.ID,
		"create_time":   shareBill.CreatedAt,
	})

}

/**

  1. 校验拼单号 地址
  2. 校验商品信息
  3. 参与拼团
  4. 创建订单
  5. 更新商品
  6. 检查人数，更新拼团和订单


*/

// ShareBillJoinHandler 参与拼团
func ShareBillJoinHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	shareBillJoinModel := &ShareBillJoinModel{}
	if !utils.QuickBind(c, shareBillJoinModel) {
		return
	}

	// 1. 校验拼单号 地址
	shareBill, msgCode, _ := dbop.ShareBillCheck(&model.ShareBill{
		ID:      shareBillJoinModel.ShareBillID,
		OwnerID: shareBillJoinModel.OwnerID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUnMatchedID(c)
		return
	}

	if shareBill[0].Status == model.ShareBillFailed || shareBill[0].Status == model.ShareBillSuccess {
		code.GinShareBillDone(c)
		return
	}

	members, msgcode12, _ := dbop.ShareBillTeamCheck(&model.ShareBillTeam{ShareBillID: shareBillJoinModel.ShareBillID})

	if msgcode12.Code == code.CheckError || msgcode12.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	// 不能重复拼团
	for _, member := range members {
		if member.MemberID == ID {
			code.GinUserInTeam(c)
			return
		}
	}

	// 校验商品数量
	commodity, err := beforeCreateCheckCommodity(c, &model.CommodityInfo{ID: shareBill[0].CommodityID})
	if !err {
		return
	}

	if commodity.Status != model.CommodityStatusOnShelf {
		code.GinNotOnShelf(c)
		return
	}

	// 校验地址
	_, err2 := beforeCreateCheckAddress(c, &model.CustomerAddress{
		ID:         shareBillJoinModel.AddressID,
		CustomerID: ID,
	})

	if !err2 {
		return
	}

	// 拼单中可继续
	tx := model.Db.Self.Begin()

	// 加入拼单
	_, msgCode3, _ := dbop.ShareBillTeamCreate(tx, &model.ShareBillTeam{
		ShareBillID: shareBill[0].ID,
		MemberID:    ID,
	})

	if msgCode3.Code == code.InsertError || msgCode3.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 创建个人订单
	_, msgCode5, _ := dbop.OrderCreate(tx, &model.Order{
		ID:          utils.OrderIDGenerate(),
		ShareBillID: shareBill[0].ID,
		CustomerID:  ID,
		AddressID:   shareBillJoinModel.AddressID,
		DueAt:       time.Now(),
		CommodityAt: time.Now(),
		FinishAt:    time.Now(),
		Status:      model.OrderCreated,
	})

	if msgCode5.Code == code.InsertError || msgCode5.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 更新商品
	_, msgCode6, _ := dbop.CommodityInfoUpdate(tx, &model.CommodityInfo{
		ID: shareBill[0].CommodityID,
	}, &model.CommodityInfo{
		Count:  commodity.Count - 1,
		Status: CommodityStatusDecide(commodity.Status, true, commodity.Count-1),
	})

	if msgCode6.Code == code.UpdateError {
		tx.Rollback()
		code.GinServerError(c)
		return
	} else if msgCode6.Code == code.DBEmpty {
		tx.Rollback()
		code.GinUnMatchedID(c)
		return
	}

	// 更新拼单订单
	if shareBill[0].Status == model.ShareBillWaitingForTwo {
		// 更新拼单订单
		_, msgCode7, _ := dbop.ShareBillUpdate(tx,
			&model.ShareBill{ID: shareBill[0].ID},
			&model.ShareBill{Status: model.ShareBillWaitingForOne})

		if msgCode7.Code == code.UpdateError || msgCode7.Code == code.DBEmpty {
			tx.Rollback()
			code.GinServerError(c)
			return
		}

	} else {
		// 5.更新拼单订单
		_, msgCode7, _ := dbop.ShareBillUpdate(tx,
			&model.ShareBill{ID: shareBill[0].ID},
			&model.ShareBill{Status: model.ShareBillSuccess, FinishAt: time.Now()})

		if msgCode7.Code == code.UpdateError || msgCode7.Code == code.DBEmpty {
			tx.Rollback()
			code.GinServerError(c)
			return
		}

		// 5. 更新订单
		// get
		orders, msgCode8, _ := dbop.OrderCheck(&model.Order{
			ShareBillID: shareBill[0].ID,
		})

		if msgCode8.Code == code.CheckError || msgCode8.Code == code.DBEmpty {
			tx.Rollback()
			code.GinServerError(c)
			return
		}
		// update
		for _, value := range orders {
			_, msgCode9, _ := dbop.OrderUpdate(tx, &model.Order{ID: value.ID}, &model.Order{
				Status: model.OrderDue,
				DueAt:  time.Now()})

			if msgCode9.Code == code.UpdateError {
				tx.Rollback()
				code.GinServerError(c)
				return
			}
		}
	}

	tx.Commit()
	code.GinOKEmpty(c)
}

// ShareBillLink 获取拼单链接
//func ShareBillLink(c *gin.Context) {}

// ShareBillListHandler 获取拼单列表
func ShareBillListHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	limit := c.Query("limit")
	page := c.Query("page")

	shareBills, msgCode, _ := dbop.ShareBillLimitPageCheck(&model.ShareBill{
		OwnerID: ID,
	}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinOKPayload(c, &gin.H{
			"list": []*ShareBillDetailsModel{},
		})
		return
	}

	var shareBillDetailModels []*ShareBillDetailsModel = make([]*ShareBillDetailsModel, len(shareBills))

	for index, shareBill := range shareBills {

		shareBillDetailSingle, msgCode2, _ := shareBillDetailSingleUnion(shareBill)

		if msgCode2.Code == code.CheckError || msgCode2.Code == code.DBEmpty {
			code.GinServerError(c)
			return
		}

		shareBillDetailModels[index] = shareBillDetailSingle

	}
	code.GinOKPayload(c, &gin.H{
		"list":  shareBillDetailModels,
		"count": len(shareBillDetailModels),
	})
}

// ShareBillDetailHandler 获取拼单订单细节
func ShareBillDetailHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	if !utils.QuickBindPath(c, pathStringIDModel) {
		return
	}

	// 检查是否是成员
	_, msgCode0, _ := dbop.ShareBillTeamCheck(&model.ShareBillTeam{
		ShareBillID: pathStringIDModel.ID,
		MemberID:    ID,
	})

	if msgCode0.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode0.Code == code.DBEmpty {
		code.GinUnAuthorized(c)
		return
	}

	// 是成员
	shareBills, msgCode, _ := dbop.ShareBillCheck(&model.ShareBill{
		ID: pathStringIDModel.ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinOKPayloadAny(c, &ShareBillDetailsModel{})
		return
	}

	shareBillDetailModel, msgCode2, _ := shareBillDetailSingleUnion(shareBills[0])

	if msgCode2.Code == code.CheckError || msgCode2.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKPayloadAny(c, shareBillDetailModel)
}

func shareBillDetailSingleUnion(shareBill *model.ShareBill) (*ShareBillDetailsModel, *code.MsgCode, bool) {
	commodity, msgCode2, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{ID: shareBill.CommodityID})
	if msgCode2.Code == code.CheckError || msgCode2.Code == code.DBEmpty {
		return nil, msgCode2, false
	}

	members, msgCode3, _ := dbop.ShareBillTeamUnionCheck(&model.ShareBillTeam{ShareBillID: shareBill.ID})
	if msgCode3.Code == code.CheckError || msgCode3.Code == code.DBEmpty {
		return nil, msgCode3, false
	}

	merchant, msgCode4, _ := dbop.MerchantInfoCheck(&model.MerchantInfo{MerchantID: commodity[0].MerchantID})
	if msgCode4.Code == code.CheckError || msgCode4.Code == code.DBEmpty {
		return nil, msgCode4, false
	}

	shareBillDetailModels := &ShareBillDetailsModel{
		CommodityInfo: commodity[0],
		CreateAt:      shareBill.CreatedAt,
		DoneAt:        shareBill.FinishAt,
		Member:        members,
		OwnerID:       shareBill.OwnerID,
		ShareBillID:   shareBill.ID,
		ShopName:      merchant.ShopName,
		ShopAvatar:    merchant.ShopAvatar,
		Status:        shareBill.Status,
	}

	return shareBillDetailModels, &code.MsgCode{Msg: "OK", Code: code.OK}, true
}
