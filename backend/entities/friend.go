package entities

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	FromUserID   string `gorm:"not null" json:"fromUserID"`
	FromUser     User   `gorm:"foreignKey:FromUserID"`
	ToUserID     string `gorm:"not null" json:"toUserID"`
	ToUser       User   `gorm:"foreignKey:ToUserID"`
	FriendStatus string `gorm:"not null" json:"friendStatus"`
}
