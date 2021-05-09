package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username      string     `json:"username"`
	Email         string     `gorm:"type:varchar(100);unique_index"`
	Password      string     `json:"password"`
	Role          string     `json:"role"`
	EmailVerified bool       `gorm:"DEFAULT:false" json:"emailVerified"`
	Overlays      []*Overlay `gorm:"many2many:user_overlays;"`
}

type UserResolver struct {
	U User
}

func (u *UserResolver) Username() string {
	return u.U.Username
}

func (u *UserResolver) Email() string {
	return u.U.Email
}

func (u *UserResolver) Role() string {
	return u.U.Role
}

func (u *UserResolver) EmailVerified() bool {
	return u.U.EmailVerified
}
