package model

type UserAuthor struct {
	ID         uint   `gorm:"<-:create;not null;uniqueIndex"`
	Phone      string `gorm:"<-;not null;unique;type:char(16)"`
	Type       uint8  `gorm:"<-:create;not null"`
	Salt       string `gorm:"<-:create;not null;type:char(32)"`
	Password   string `gorm:"<-;not null;type:char(32)"`
	PrivateKey string `gorm:"<-;not null;type:char(128)"`
}

//type AdministratorAuthor struct {
//	AdminAuthor UserAuthor `gorm:"embedded"`
//	Name        string     `gorm:"<-;not null;type:char(32)"`
//}
