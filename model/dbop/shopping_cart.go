package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
	"strconv"
)

type ShoppingCartUnionModel struct {
	CartID     uint    `json:"cart_id"`
	Count      uint    `json:"count"` // 商品库存
	ID         uint    `json:"id"`
	Intro      string  `json:"intro"`
	MerchantID uint    `json:"merchant_id"`
	Name       string  `json:"name"`
	Picture    string  `json:"picture"`
	Price      float64 `json:"price"` // 商品单价
	Status     uint    `json:"status"`
}

// ShoppingCartCreate InsertError
func ShoppingCartCreate(tx *gorm.DB, cart *model.ShoppingCart) (*model.ShoppingCart, *code.MsgCode, error) {

	// 数据库存储
	result := tx.Create(cart)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return cart, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func ShoppingCartDrop(tx *gorm.DB, condition *model.ShoppingCart) (*code.MsgCode, error) {

	result := tx.Delete(condition)

	if result.Error != nil {
		return &code.MsgCode{Msg: "DropError", Code: code.DropError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func ShoppingCartCheck(condition *model.ShoppingCart) ([]*model.ShoppingCart, *code.MsgCode, error) {

	var searchShoppingCart []*model.ShoppingCart

	// 条件由外部决定
	result := model.Db.Self.Where(condition).
		Not(&model.ShoppingCart{}).
		Find(&searchShoppingCart)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchShoppingCart, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func ShoppingCartLimitPageUnionCheck(condition *model.ShoppingCart, limit, page string) ([]*ShoppingCartUnionModel, *code.MsgCode, error) {

	var shoppingCartUnionModel []*ShoppingCartUnionModel

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	var result *gorm.DB

	if limitInt == 0 {
		limitInt = -1
		pageInt = 1
	}

	result = model.Db.Self.
		Table("shopping_carts").
		Select(
			"shopping_carts.id as cart_id," +
				"commodity_infos.count as count," +
				"commodity_infos.id as id," +
				"commodity_infos.intro as intro," +
				"commodity_infos.merchant_id as merchant_id," +
				"commodity_infos.name as name," +
				"commodity_infos.picture as picture," +
				"commodity_infos.price as price," +
				"commodity_infos.status as status").
		Joins("left join commodity_infos on `commodity_infos`.`id` = `shopping_carts`.`commodity_id`").
		Where(condition).
		Offset(limitInt * pageInt).
		Limit(limitInt).
		Find(&shoppingCartUnionModel)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return shoppingCartUnionModel, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
