package dbop

import (
    "HappyShopTogether/model"
    "HappyShopTogether/utils/code"
    "gorm.io/gorm"
    "strconv"
)

// ShareBillCreate InsertError
func ShareBillCreate(tx *gorm.DB, shareBill *model.ShareBill) (*model.ShareBill, *code.MsgCode, error) {

	// 数据库存储
	result := tx.Create(shareBill)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return shareBill, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

//
//
//func ShareBillDrop(tx *gorm.DB, condition *model.ShareBill) (*code.MsgCode, error) {
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

func ShareBillCheck(condition *model.ShareBill) ([]*model.ShareBill, *code.MsgCode, error) {

    var searchShareBill []*model.ShareBill

    // 条件由外部决定
    result := model.Db.Self.Where(condition).
        Not(&model.ShareBill{}).
        Find(&searchShareBill)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchShareBill, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func ShareBillLimitPageCheck(condition *model.ShareBill, limit, page string) ([]*model.ShareBill, *code.MsgCode, error) {

    var searchShareBill []*model.ShareBill

    limitInt, _ := strconv.Atoi(limit)
    pageInt, _ := strconv.Atoi(page)

    if limitInt == 0 || pageInt == 0 {
        return ShareBillCheck(condition)
    }

    // 条件由外部决定
    result := model.Db.Self.
        Where(condition).
        Not(&model.ShareBill{}).
        Limit(limitInt).
        Offset(limitInt * pageInt).
        Find(&searchShareBill)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchShareBill, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


// ShareBillUpdate  shareBill 就是要保存的数据
func ShareBillUpdate(tx *gorm.DB, condition, shareBill *model.ShareBill) (*model.ShareBill, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(shareBill)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
