package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
)



func addressDecide(current, update uint8) uint8 {
	if current == model.AddressDelete {
		return current
	}
	return update
}

// CustomerAddressCreate InvalidPhone InsertError OK
func CustomerAddressCreate(tx *gorm.DB, address *model.CustomerAddress) (*model.CustomerAddress, *code.MsgCode, error) {

	if !utils.CheckMobile(address.Phone) {
		return nil, &code.MsgCode{Msg: "InvalidPhone", Code: code.InvalidPhone}, nil
	}

	_, msgCode, _ := CustomerAddressCheck(&model.CustomerAddress{
		CustomerID: address.CustomerID,
	})

	if msgCode.Code == code.CheckError {
		return nil, msgCode, nil
	} else if msgCode.Code == code.DBEmpty {
		address.Default = addressDecide(model.AddressNotExsit, model.AddressDefault)
	} else {
		address.Default = addressDecide(model.AddressNotExsit, model.AddressNotDefault)
	}

	result := tx.Create(address)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	return address, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
func CustomerAddressDrop(tx *gorm.DB, address *model.CustomerAddress) (*code.MsgCode, error) {

	result := tx.Model(address).Updates(&model.CustomerAddress{
		Default: addressDecide(model.AddressNotExsit, model.AddressDelete),
	})

	if result.Error != nil {
		return &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil

}
func CustomerAddressCheck(address *model.CustomerAddress) ([]*model.CustomerAddress, *code.MsgCode, error) {

	var addresses []*model.CustomerAddress

	// 条件由外部决定
	result := model.Db.Self.Where(address).Find(&addresses)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return addresses, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func CustomerAddressUpdate(tx *gorm.DB, condition, address *model.CustomerAddress) (*model.CustomerAddress, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(address)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
