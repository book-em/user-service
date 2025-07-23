package model

import (
	"gorm.io/gorm"
)
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null;type:string"`
	Password string `json:"password" gorm:"not null;type:string"`
	Email    string `json:"email" gorm:"unique;not null;type:string"`
}
