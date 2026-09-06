package entities

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	Model
	Name         string `gorm:"not null" json:"name"`
	UserName     string `gorm:"not null; unique" json:"userName"`
	Mobilenumber string `gorm:"not null; unique" json:"mobileNumber"`
	Password     string `gorm:"not null" json:"password"`
	PubKey       string `gorm:"not null" json:"pubKey"`
}
