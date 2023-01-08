package model
//
//import "time"
//
//type ContactList struct {
//	SelfID        uint `gorm:"primaryKey;autoIncrement:false;<-:create;not null;Index"`
//	Self          UserAuthor
//	ContactWithID uint `gorm:"primaryKey;autoIncrement:false;<-:create;not null"`
//	ContactWith   UserAuthor
//	UnreadMsg     uint      `gorm:"<-;not null"`
//	LastContact   time.Time `gorm:"<-;not null"`
//}
//
//type ContactMsg struct {
//	ID      uint `gorm:"<-:create;not null;uniqueIndex"`
//	FromID  uint `gorm:"<-:create;not null"`
//	From    UserAuthor
//	ToID    uint       `gorm:"<-:create;not null"`
//	To      UserAuthor
//	SendAt  time.Time  `gorm:"<-:create;not null"`
//	status  bool       `gorm:"<-;not null"`
//	Content string     `gorm:"<-:create;not null;type:varchar(512)"`
//}
//
//type DailyChatRecord struct {
//	ID         uint `gorm:"<-:create;not null:uniqueIndex"`
//	OneID      uint `gorm:"<-:create;not null"`
//	One        UserAuthor
//	TheOtherID uint `gorm:"<-:create;not null"`
//	TheOther   UserAuthor
//	RecordDay  time.Time `gorm:"<-:create;not null"`
//	RecordNo   uint      `gorm:"<-:create;not null"`
//	Content    string    `gorm:"<-;not null"`
//}
