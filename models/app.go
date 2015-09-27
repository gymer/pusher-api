package models

import "github.com/jinzhu/gorm"

type App struct {
	gorm.Model
	Name   string
	UserID uint `sql:"index"`
}
