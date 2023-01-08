package dbop

import (
    "HappyShopTogether/model"
    "HappyShopTogether/utils/code"
    "gorm.io/gorm"
    "strconv"
)

// OrderCreate InsertError
func OrderCreate(tx *gorm.DB, order *model.Order) (*model.Order, *code.MsgCode, error) {

	// 数据库存储
    result := tx.Create(order)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
    return order, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}




//func OrderDrop(tx *gorm.DB, condition *model.Order) (*code.MsgCode, error) {
//
//    result := tx.Delete(condition)
//
//    if result.Error != nil {
//        return &code.MsgCode{Msg: "DropError", Code: code.DropError}, result.Error
//    }
//
//    // 找不到用户
//    if result.RowsAffected == 0 {
//        return &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
//    }
//
//    return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
//}


func OrderCheck(condition *model.Order) ([]*model.Order, *code.MsgCode, error) {

    var searchOrder []*model.Order

    // 条件由外部决定
    result := model.Db.Self.Where(condition).
        Not(&model.Order{}).
        Find(&searchOrder)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchOrder, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func OrderLimitPageCheck(condition *model.Order, limit, page string) ([]*model.Order, *code.MsgCode, error) {

    var searchOrder []*model.Order

    limitInt, _ := strconv.Atoi(limit)
    pageInt, _ := strconv.Atoi(page)

    if limitInt == 0 || pageInt == 0 {
        return OrderCheck(condition)
    }

    // 条件由外部决定
    result := model.Db.Self.
        Where(condition).
        Not(&model.Order{}).
        Limit(limitInt).
        Offset(limitInt * pageInt).
        Find(searchOrder)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchOrder, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


// OrderUpdate  order 就是要保存的数据
func OrderUpdate(tx *gorm.DB, condition, order *model.Order) (*model.Order, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(order)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
