package handler

import (
	"HappyShopTogether/model"
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
)

type AddressModel struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

type AddressUpdateModel struct {
	Address string `json:"address"`
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}

type AddressCheckModel struct {
    Address string `json:"address"`
    ID      uint   `json:"id"`
    Name    string `json:"name"`
    Phone   string `json:"phone"`
    Default uint8 `json:"default"`
}

type AddressPathModel struct {
	AddressID uint `uri:"addressID" binding:"required"`
}

func AddressAddHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	addressToAdd := &AddressModel{}
	utils.QuickBind(c, addressToAdd)

	_, msgCode, _ := dbop.CustomerAddressCreate(model.Db.Self, &model.CustomerAddress{
		CustomerID:   ID,
		Address:      addressToAdd.Address,
		Phone:        addressToAdd.Phone,
		ReceiverName: addressToAdd.Name,
	})

	if msgCode.Code == code.InsertError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.InvalidPhone {
		code.GinBadRequest(c)
		return
	}

	code.GinOKEmpty(c)
}
func AddressUpdateHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	addressUpdateModel := &AddressUpdateModel{}
	utils.QuickBind(c, addressUpdateModel)

	addresses, msgCode, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{
		ID: addressUpdateModel.ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	if addresses[0].CustomerID != ID {
		code.GinUnAuthorized(c)
		return
	}

	_, msgCode2, _ := dbop.CustomerAddressUpdate(
		model.Db.Self,
		&model.CustomerAddress{
			ID: addressUpdateModel.ID,
		},
		&model.CustomerAddress{
			Address:      addressUpdateModel.Address,
			Phone:        addressUpdateModel.Phone,
			ReceiverName: addressUpdateModel.Name,
		},
	)

	if msgCode2.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
func AddressDefaultHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	addressPathModel := &AddressPathModel{}
	utils.QuickBindPath(c, addressPathModel)

	// 查出原来的default
	defaultAddress, msgCode, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{
		CustomerID: ID,
		Default:    model.AddressDefault,
	})

	flag := true

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		flag = false
	}

	tx := model.Db.Self.Begin()

	// 0 -> 1

	_, msgCode2, _ := dbop.CustomerAddressUpdate(
		tx,
		&model.CustomerAddress{
			ID:         addressPathModel.AddressID,
			CustomerID: ID,
		},
		&model.CustomerAddress{
			Default: model.AddressDefault,
		},
	)

	if msgCode2.Code == code.UpdateError {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	// 1 -> 0
	if flag {
		_, msgCode3, _ := dbop.CustomerAddressUpdate(
			tx,
			&model.CustomerAddress{
				ID:         defaultAddress[0].ID,
				CustomerID: ID,
			},
			&model.CustomerAddress{
				Default: model.AddressNotDefault,
			},
		)

		if msgCode3.Code == code.UpdateError {
			tx.Rollback()
			code.GinServerError(c)
			return
		}
	}

	tx.Commit()

	code.GinOKEmpty(c)
}

func AddressCheckHander(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	addresses, msgCode, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{
		CustomerID: ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

    var ades = make([]*AddressCheckModel, len(addresses))

	for index, address := range addresses {
		ades[index] = &AddressCheckModel{
			Address: address.Address,
			ID:      address.ID,
			Phone:   address.Phone,
			Name:    address.ReceiverName,
            Default: address.Default,
		}
	}

	code.GinOKPayload(c, &gin.H{
		"address": ades,
        "count": len(ades),
    })
}

func AddressDeleteHandler(c *gin.Context) {
	ID, _, _ := utils.GetTokenInfo(c)

	addressPathModel := &AddressPathModel{}
	utils.QuickBindPath(c, addressPathModel)

	msgCode, _ := dbop.CustomerAddressDrop(model.Db.Self, &model.CustomerAddress{
		ID:         addressPathModel.AddressID,
		CustomerID: ID,
	})

	if msgCode.Code == code.DropError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinOKEmpty(c)
		return
	}

	// 删除之后找一条 改成 default  性能原因直接用数据库

	first := &model.CustomerAddress{
		CustomerID: ID,
	}

	result := model.Db.Self.Limit(1).Find(first)

	if result.Error != nil {
		code.GinServerError(c)
		return
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		code.GinOKEmpty(c)
		return
	}

	_, msgCode2, _ := dbop.CustomerAddressUpdate(
		model.Db.Self,
		&model.CustomerAddress{
			ID: first.ID,
		},
		&model.CustomerAddress{
			Default: model.AddressDefault,
		},
	)

	if msgCode2.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
