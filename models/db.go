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
	DBconf, err := config.NewConfig("ini", "config/database.conf")

	if err != nil {
		panic(err)
	}

	config, err := DBconf.GetSection(env)

	if err != nil {
		panic(err)
	}

	DB, err = gorm.Open("postgres", "host="+config["host"]+" user="+config["user"]+" dbname="+config["dbname"]+" sslmode=disable")

	if err != nil {
		log.Fatalln(err)
		return
	}

	switch env {
	case "dev":
		// DB.LogMode(true)
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
