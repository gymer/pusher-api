package models

import (
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db  gorm.DB
	err error
)

func init() {
	db, err = gorm.Open("postgres", "user=postgres dbname=gymer_dev sslmode=disable")

	if err != nil {
		beego.Error(err)
		return
	}

	db.DB()
	migrations()

	if beego.RunMode == "dev" {
		seed()
	}
}

func migrations() {
	db.AutoMigrate(&User{}, &App{})
}

func seed() {
	user := User{}
	app := App{}
	db.FirstOrCreate(&user, User{Email: "test@test.com", Password: "qwerty"})
	db.FirstOrCreate(&app, App{Name: "Test App", UserID: user.ID})
}
