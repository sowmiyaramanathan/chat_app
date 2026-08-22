package entities

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	FromUserID   uint64 `gorm:"not null" json:"fromUserID"`
	FromUser     User   `gorm:"foreignKey:FromUserID"`
	ToUserID     uint64 `gorm:"not null" json:"toUserID"`
	ToUser       User   `gorm:"foreignKey:ToUserID"`
	FriendStatus string `gorm:"not null" json:"friendStatus"`
}
