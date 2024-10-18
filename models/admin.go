package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	AdminName string
	Account   string
	Password  string
}
