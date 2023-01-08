package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
)

/*
type CustomerInfo struct {
    CustomerID uint `gorm:"primaryKey;<-:create;not null;uniqueIndex"`
    Customer   UserAuthor
    NickName   string    `gorm:"<-;not null;type:char(32)"`
    Avatar     string    `gorm:"<-;type:varchar(256)"`
    Name       string    `gorm:"<-:create;not null;type:char(32)"`
    Birth      time.Time `gorm:"<-:create;not null"`
    Intro      string    `gorm:"<-;type:varchar(256)"`
}
*/

// CustomerInfoCreate @todo
func CustomerInfoCreate(tx *gorm.DB, customer *model.CustomerInfo) (*model.CustomerInfo, *code.MsgCode, error) {

	// 重复校验
	_, msgCode, err := CustomerInfoCheck(&model.CustomerInfo{
		CustomerID: customer.CustomerID,
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
	result := tx.Create(customer)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return customer, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

//func CustomerInfoDrop(customer *model.CustomerInfo) (*model.CustomerInfo, *code.MsgCode, error) {
//	return nil, nil, nil
//}

// CustomerInfoCheck CheckError DbEmpty OK
func CustomerInfoCheck(customer *model.CustomerInfo) (*model.CustomerInfo, *code.MsgCode, error) {
	searchCustomer := &model.CustomerInfo{}

	// 条件由外部决定
	result := model.Db.Self.Where(customer).Find(searchCustomer)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchCustomer, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// CustomerInfoUpdate UpdateError DBEmpty OK
func CustomerInfoUpdate(tx *gorm.DB, condition, customer *model.CustomerInfo) (*model.CustomerInfo, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(customer)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
