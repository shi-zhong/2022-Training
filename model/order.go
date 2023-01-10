package model

import "time"

type Order struct {
	ID          string `gorm:"<-:create;type:char(24);not null;uniqueIndex"`
	ShareBillID string `gorm:"<-:create;not null"`
	ShareBill   ShareBill
	CustomerID  uint `gorm:"<-:create;not null"`
	Customer    UserAuthor
    MerchantID  uint `gorm:"<-:create;not null"`
    Merchant    UserAuthor
	AddressID   uint `gorm:"<-:create;not null"`
	Address     CustomerAddress
	Status      uint      `gorm:"<-;not null"`
	CreatedAt   time.Time `gorm:"<-:create;not null"`
	DueAt       time.Time `gorm:"<-:create"`
	CommodityAt time.Time `gorm:"<-:create"`
	FinishAt    time.Time `gorm:"<-:create"`
}
