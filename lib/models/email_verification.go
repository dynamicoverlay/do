package models

import "github.com/jinzhu/gorm"

type EmailVerification struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index"`
	Code  string `gorm:"type:varchar(255)"`
}
