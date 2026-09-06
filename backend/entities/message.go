package entities

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromUserID string `gorm:"not null,index:idx_messages_users" json:"fromUserID"`
	FromUser   User   `gorm:"foreignkey:FromUserID"`
	ToUserID   string `gorm:"not null,index:idx_messages_users" json:"toUserID"`
	ToUser     User   `gorm:"foreignkey:ToUserID"`
	Message    string `gorm:"not null" json:"message"`
}
