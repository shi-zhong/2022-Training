package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
	"time"
)

type orderUnionModel struct {
	Order     *model.Order
	Commodity *model.CommodityInfo
	Merchant  *model.MerchantInfo
	Address   *model.CustomerAddress
}

func orderUnionCheck(c *gin.Context, order *model.Order) (*orderUnionModel, bool) {

	// 为了查 商品id
	sharebill, msgCode2, _ := dbop.ShareBillCheck(&model.ShareBill{
		ID: order.ShareBillID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingShareBill(c)
		return nil, false
	}

	commodity, msgCode3, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		ID: sharebill[0].CommodityID,
	})

	if msgCode3.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode3.Code == code.DBEmpty {
		code.GinMissingShareBill(c)
		return nil, false
	}

	merchant, msgCode4, _ := dbop.MerchantInfoCheck(&model.MerchantInfo{
		MerchantID: commodity[0].MerchantID,
	})

	if msgCode4.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode4.Code == code.DBEmpty {
		code.GinMissingMerchant(c)
		return nil, false
	}

	address, msgCode5, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{ID: order.AddressID})

	if msgCode5.Code == code.CheckError {
		code.GinServerError(c)
		return nil, false
	} else if msgCode5.Code == code.DBEmpty {
		code.GinMissingMerchant(c)
		return nil, false
	}

	return &orderUnionModel{order, commodity[0], merchant, address[0]}, true
}

func OrderCustomerListHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	limit := c.Query("limit")
	page := c.Query("page")

	orders, msgCode, _ := dbop.OrderLimitPageCheck(&model.Order{CustomerID: ID}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	var orderDetails []*orderDetailsModel = make([]*orderDetailsModel, len(orders))

	for index, order := range orders {
		orderUnion, err := orderUnionCheck(c, order)
		if !err {
			return
		}
		orderdetail := orderDetailSingleUnion(orderUnion)
		orderDetails[index] = orderdetail
	}

	code.GinOKPayload(c, &gin.H{
		"list":  orderDetails,
		"count": len(orderDetails),
	})
}

func OrderDetailHandler(c *gin.Context) {
	ID, Type, _ := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode, _ := dbop.OrderCheck(&model.Order{ID: pathStringIDModel.ID})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	// order commodity merchant address
	orderUnion, err2 := orderUnionCheck(c, order[0])
	if !err2 {
		return
	}

	// 不是该订单的商家或者顾客
	if (Type == model.UserTypeCustomer && orderUnion.Order.CustomerID != ID) ||
		(Type == model.UserTypeMerchant && orderUnion.Commodity.MerchantID != ID) {
		code.GinUnAuthorized(c)
		return
	}

	// 组合

	orderDetails := orderDetailSingleUnion(orderUnion)

	code.GinOKPayloadAny(c, orderDetails)
}

func OrderMerchantConfirmHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode, _ := dbop.OrderCheck(&model.Order{ID: pathStringIDModel.ID})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	// order commodity merchant address
	orderUnion, err2 := orderUnionCheck(c, order[0])
	if !err2 {
		return
	}

	// 不是该订单的商家
	if orderUnion.Commodity.MerchantID != ID {
		code.GinUnAuthorized(c)
		return
		// 订单未下单
	} else if orderUnion.Order.Status != model.OrderDue {
		code.GinOrderNotInDue(c)
		return
	}

	_, msgCode2, _ := dbop.OrderUpdate(model.Db.Self,
		&model.Order{ID: pathStringIDModel.ID},
		&model.Order{Status: model.OrderCommodity, CommodityAt: time.Now()})

	if msgCode2.Code == code.UpdateError || msgCode2.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

func OrderCustomerConfirmHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode, _ := dbop.OrderCheck(&model.Order{ID: pathStringIDModel.ID, CustomerID: ID})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if order[0].Status != model.OrderCommodity {
		code.GinOrderNotInCommodity(c)
		return
	}

	_, msgCode2, _ := dbop.OrderUpdate(model.Db.Self,
		&model.Order{ID: pathStringIDModel.ID, CustomerID: ID},
		&model.Order{Status: model.OrderFinish, FinishAt: time.Now()})

	if msgCode2.Code == code.UpdateError || msgCode2.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

// 订单有关时间
type orderAtModel struct {
	CommodityAt time.Time `json:"commodity_at"` // 订单发货时间
	CreateAt    time.Time `json:"create_at"`    // 订单创建时间(付款时间)
	DueAt       time.Time `json:"due_at"`       // 订单下单时间(拼单完成时间)
	FinishAt    time.Time `json:"finish_at"`    // 订单结束时间(收货/失败时间)
}

type addressModel struct {
	ID      uint   `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

// order
type orderDetailsModel struct {
	Address       *addressModel        `json:"address,omitempty"`
	CommodityInfo *model.CommodityInfo `json:"commodity_info"`
	MerchantID    uint                 `json:"merchant_id"`
	OrderAt       *orderAtModel        `json:"order_at"` // 订单有关时间
	OrderID       string               `json:"order_id"`
	ShareBillID   string               `json:"share_bill_id"`
	ShopName      string               `json:"shop_name"`
	Status        uint                 `json:"status"`
}

func orderDetailSingleUnion(unionModel *orderUnionModel) *orderDetailsModel {

	orderAddress := &addressModel{
		ID:      unionModel.Address.ID,
		Address: unionModel.Address.Address,
		Name:    unionModel.Address.ReceiverName,
		Phone:   unionModel.Address.Phone,
	}

	orderAt := &orderAtModel{
		CommodityAt: unionModel.Order.CommodityAt, // 订单发货时间
		CreateAt:    unionModel.Order.CreatedAt,   // 订单创建时间(付款时间)
		DueAt:       unionModel.Order.DueAt,       // 订单下单时间(拼单完成时间)
		FinishAt:    unionModel.Order.FinishAt,    // 订单结束时间(收货/失败时间)
	}

	return &orderDetailsModel{
		Address:       orderAddress,
		CommodityInfo: unionModel.Commodity,
		MerchantID:    unionModel.Merchant.MerchantID,
		OrderAt:       orderAt,
		OrderID:       unionModel.Order.ID,
		ShareBillID:   unionModel.Order.ShareBillID,
		ShopName:      unionModel.Merchant.ShopName,
		Status:        unionModel.Order.Status,
	}
}
