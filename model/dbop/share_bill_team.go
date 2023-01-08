package dbop

import (
	"HappyShopTogether/model"
	"HappyShopTogether/utils/code"
	"gorm.io/gorm"
//	"strconv"
	"time"
)

type Member struct {
    Avatar   string    `json:"avatar"`
    ID       uint      `json:"id"`
    JoinTime time.Time `json:"join_time"`
    Nickname string    `json:"nickname"`
}

// ShareBillTeamCreate InsertError
func ShareBillTeamCreate(tx *gorm.DB, shareBillTeam *model.ShareBillTeam) (*model.ShareBillTeam, *code.MsgCode, error) {

	// 数据库存储
	result := tx.Create(shareBillTeam)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return shareBillTeam, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
//
//func ShareBillTeamDrop(tx *gorm.DB, condition *model.ShareBillTeam) (*code.MsgCode, error) {
//
//	result := tx.Delete(condition)
//
//	if result.Error != nil {
//		return &code.MsgCode{Msg: "DropError", Code: code.DropError}, result.Error
//	}
//
//	// 找不到用户
//	if result.RowsAffected == 0 {
//		return &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
//	}
//
//	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
//}

func ShareBillTeamCheck(condition *model.ShareBillTeam) ([]*model.ShareBillTeam, *code.MsgCode, error) {

	var searchShareBillTeam []*model.ShareBillTeam

	// 条件由外部决定
	result := model.Db.Self.Where(condition).
		Not(&model.ShareBillTeam{}).
		Find(&searchShareBillTeam)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchShareBillTeam, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
//
//func ShareBillTeamLimitPageCheck(condition *model.ShareBillTeam, limit, page string) ([]*model.ShareBillTeam, *code.MsgCode, error) {
//
//	var searchShareBillTeam []*model.ShareBillTeam
//
//	limitInt, _ := strconv.Atoi(limit)
//	pageInt, _ := strconv.Atoi(page)
//
//	if limitInt == 0 || pageInt == 0 {
//		return ShareBillTeamCheck(condition)
//	}
//
//	// 条件由外部决定
//	result := model.Db.Self.
//		Where(condition).
//		Not(&model.ShareBillTeam{}).
//		Limit(limitInt).
//		Offset(limitInt * pageInt).
//		Find(searchShareBillTeam)
//
//	if result.Error != nil {
//		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
//	}
//
//	// 找不到用户
//	if result.RowsAffected == 0 {
//		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
//	}
//
//	return searchShareBillTeam, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
//}


func ShareBillTeamUnionCheck(condition *model.ShareBillTeam) ([]*Member, *code.MsgCode, error) {

	var members []*Member

	var result *gorm.DB

	result = model.Db.Self.
		Table("share_bill_teams").
		Select(
			"share_bill_teams.created_at as join_time," +
                "customer_infos.customer_id as id," +
				"customer_infos.avatar as avatar," +
				"customer_infos.nick_name as nickname",
		).
        Joins("join customer_infos on `customer_infos`.`customer_id` = `share_bill_teams`.`member_id`").
		Where(condition).
		Find(&members)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return members, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// ShareBillTeamUpdate  shareBillTeam 就是要保存的数据
//func ShareBillTeamUpdate(tx *gorm.DB, condition, shareBillTeam *model.ShareBillTeam) (*model.ShareBillTeam, *code.MsgCode, error) {
//	result := tx.Model(condition).Updates(shareBillTeam)
//
//	if result.Error != nil {
//		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
//	}
//
//	if result.RowsAffected == 0 {
//		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
//	}
//
//	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
//}
