package models

import (
	"fmt"
	"log"

	"gopkg.in/testfixtures.v1"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	DB  gorm.DB
	err error
)

func ConnectDB() {
	DBconf, err := config.NewConfig("ini", "conf/database.conf")

	if err != nil {
		panic(err)
	}

	var host = DBconf.String(beego.RunMode + "::host")
	var user = DBconf.String(beego.RunMode + "::user")
	var dbname = DBconf.String(beego.RunMode + "::dbname")

	fmt.Printf("Runmode = %+v \n", beego.RunMode)
	fmt.Printf("DB Config = %+v \n", DBconf.String(beego.RunMode+"::dbname"))

	DB, err = gorm.Open("postgres", "host="+host+" user="+user+" dbname="+dbname+" sslmode=disable")

	if err != nil {
		beego.Error(err)
		return
	}

	DB.DB()

	if beego.RunMode == "test" {
		loadFixtures()
	}
	// migrations()
}

func loadFixtures() {
	err := testfixtures.LoadFixtures("fixtures", DB.DB(), &testfixtures.PostgreSQLHelper{})

	if err != nil {
		log.Fatal(err)
	}
}

// func seed() {
// 	user := User{}
// 	app := App{}
// 	DB.FirstOrCreate(&user, User{Email: "test@test.com", Password: "qwerty"})
// 	DB.FirstOrCreate(&app, App{Name: "Test App", UserID: user.ID})
// }
