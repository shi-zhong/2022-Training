package model

import "time"



type CustomerInfo struct {
	CustomerID uint `gorm:"primaryKey;<-:create;not null;uniqueIndex"`
	Customer   UserAuthor
	NickName   string    `gorm:"<-;not null;type:char(32)"`
	Avatar     string    `gorm:"<-;type:varchar(256)"`
	Name       string    `gorm:"<-:create;not null;type:char(32)"`
	Birth      time.Time `gorm:"<-:create;not null"`
	Intro      string    `gorm:"<-;type:varchar(256)"`
}

type CustomerAddress struct {
	ID           uint `gorm:"<-:create;not null;uniqueIndex"`
	CustomerID   uint `gorm:"<-:create;not null"`
	Customer     UserAuthor
	ReceiverName string `gorm:"<-;not null;type:char(32)"`
	Phone        string `gorm:"<-;not null;type:char(16)"`
	Address      string `gorm:"<-;not null;type:varchar(256)"`
	Default      uint8   `gorm:"<-;not null"`
}

type MerchantInfo struct {
	MerchantID uint `gorm:"primaryKey;<-:create;not null;uniqueIndex"`
	Merchant   UserAuthor
	NickName   string `gorm:"<-;not null;type:char(32)"`
	Avatar     string `gorm:"<-;type:varchar(256)"`
	Name       string `gorm:"<-:create;not null;type:char(32)"`
	Intro      string `gorm:"<-;type:varchar(256)"`
	ShopName   string `gorm:"<-;not null;unique;type:char(32)"`
	ShopIntro  string `gorm:"<-;type:varchar(256)"`
	ShopAvatar string `gorm:"<-;type:varchar(256)"`
	Address    string `gorm:"<-;not null;type:varchar(256)"`
}


