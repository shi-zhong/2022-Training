package model

//import "time"


type CommodityInfo struct {
	ID         uint `gorm:"<-:create;not null;uniqueIndex" json:"id"`
    MerchantID uint `gorm:"<-:create;not null" json:"merchant_id"`
    Merchant   UserAuthor `json:"-"`
    Count      uint    `gorm:"<-;not null" json:"count"`
    Name       string  `gorm:"<-;not null;type:char(32)" json:"name"`
    Price      float64 `gorm:"<-;not null" json:"price"`
    Intro      string  `gorm:"<-;not null;type:varchar(256)" json:"intro"`
    Status     uint    `gorm:"<-;not null" json:"status"`
    Picture    string  `gorm:"<-;not null;type:varchar(256)" json:"picture"`
}

type ShoppingCart struct {
	ID          uint `gorm:"<-:create;not null;uniqueIndex"`
	CommodityID uint `gorm:"<-:create;not null"`
	Commodity   CommodityInfo
	CustomerID  uint `gorm:"<-:create;not null"`
	Customer    UserAuthor
}
//
//type CommodityMerchantStatusLog struct {
//	ID            uint `gorm:"<-:create;not null;uniqueIndex"`
//	CommodityID   uint `gorm:"<-:create;not null"`
//	Commodity     CommodityInfo
//	OperatorID    uint `gorm:"<-:create;not null;foreignKey:ID"`
//	Operator      UserAuthor
//	OperationType uint      `gorm:"<-:create;not null"`
//	CreatedAt     time.Time `gorm:"<-:create;not null"`
//}
//
//type CommodityAdminStatusLog struct {
//	ID            uint `gorm:"<-:create;not null;uniqueIndex"`
//	CommodityID   uint `gorm:"<-:create;not null"`
//	Commodity     CommodityInfo
//	OperatorID    uint `gorm:"<-:create;not null;foreignKey:ID"`
//	Operator      AdministratorAuthor
//	OperationType uint      `gorm:"<-:create;not null"`
//	CreatedAt     time.Time `gorm:"<-:create;not null"`
//}
