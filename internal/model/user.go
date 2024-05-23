package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Nickname string
	Age      int
}

func (*User) TableName() string {
	return "user_test"
}
