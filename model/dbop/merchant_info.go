package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
)

func MerchantInfoCreate(tx *gorm.DB, merchant *model.MerchantInfo) (*model.MerchantInfo, *code.MsgCode, error) {
	// 重复校验
	_, msgCode, err := MerchantInfoCheck(&model.MerchantInfo{
		MerchantID: merchant.MerchantID,
	})

	//  数据库出错
	if msgCode.Code == code.CheckError {
		return nil, msgCode, err
	}

	if msgCode.Code == code.OK {
		// 转到更新操作
		return nil, nil, nil
	}

	// 数据库存储
	result := tx.Create(merchant)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return merchant, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

//func MerchantInfoDrop(merchant *model.MerchantInfo) (*model.MerchantInfo, *code.MsgCode, error) {
//	return nil, nil, nil
//}

func MerchantInfoCheck(merchant *model.MerchantInfo) (*model.MerchantInfo, *code.MsgCode, error) {
	searchMerchant := &model.MerchantInfo{}

	// 条件由外部决定
	result := model.Db.Self.Where(merchant).Find(searchMerchant)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到商家
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchMerchant, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// MerchantInfoUpdate UpdateError DBEmpty OK
func MerchantInfoUpdate(tx *gorm.DB, condition, merchant *model.MerchantInfo) (*model.MerchantInfo, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(merchant)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
