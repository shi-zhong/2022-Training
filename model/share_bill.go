package model

import "time"



type ShareBill struct {
	ID          string `gorm:"<-:create;type:char(24);not null;uniqueIndex"`
	OwnerID     uint   `gorm:"<-:create;not null"`
	Owner       UserAuthor
	CommodityID uint `gorm:"<-:create;not null"`
	Commodity   CommodityInfo
	CreatedAt   time.Time `gorm:"<-:create;not null"`
	FinishAt    time.Time `gorm:"<-:create"`
	Status      uint8     `gorm:"<-;not null"`
	Price       float64   `gorm:"<-:create;not null"`
}

type ShareBillTeam struct {
	ShareBillID string `gorm:"primaryKey;autoIncrement:false;<-:create;not null;Index"`
	ShareBill   ShareBill
	MemberID    uint `gorm:"primaryKey;autoIncrement:false;<-:create;not null"`
	Member      UserAuthor
	CreatedAt   time.Time `gorm:"<-:create;not null"`
}
//
//type ShareBillVisitLog struct {
//	ShareBillID string `gorm:"primaryKey;autoIncrement:false;<-:create;not null;Index"`
//	ShareBill   ShareBill
//	VisitorID   uint `gorm:"primaryKey;autoIncrement:false;<-:create;not null"`
//	Visitor     UserAuthor
//	CreatedAt   time.Time `gorm:"<-:create;not null"`
//}
