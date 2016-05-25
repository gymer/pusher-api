package models

import (
	"log"

	"gopkg.in/testfixtures.v1"

	"github.com/astaxie/beego/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	DB  gorm.DB
	err error
)

func ConnectDB(env string) {
	DBconf, err := config.NewConfig("ini", "conf/database.conf")

	if err != nil {
		panic(err)
	}

	var host = DBconf.String(env + "::host")
	var user = DBconf.String(env + "::user")
	var dbname = DBconf.String(env + "::dbname")

	DB, err = gorm.Open("postgres", "host="+host+" user="+user+" dbname="+dbname+" sslmode=disable")

	if err != nil {
		log.Fatalln(err)
		return
	}

	switch env {
	case "dev":
		DB.LogMode(true)
	case "test":
		loadFixtures()
	}
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
