package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `gorm:"type:varchar(100)" json:"first_name"`
	SecondName  string `gorm:"type:varchar(100)" json:"second_name"`
	Email       string `gorm:"unique;type:varchar(100)" json:"email"`
	PhoneNumber string `gorm:"type:varchar(15)" json:"phone_number"`
	Username    string `gorm:"unique;type:varchar(50)" json:"username"`
	Password    string `gorm:"type:varchar(255)" json:"password"`
	Role        string `gorm:"type:varchar(20)" json:"role"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
