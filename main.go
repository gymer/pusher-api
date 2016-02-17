package main

import (
	"flag"

	"github.com/gymer/pusher-api/controllers"
	_ "github.com/gymer/pusher-api/docs"
	"github.com/gymer/pusher-api/models"
	"github.com/gymer/pusher-api/routers"

	"github.com/astaxie/beego"
)

var env string
var port string

var err error

func init() {
	flag.StringVar(&env, "env", "dev", "set app environment")
	flag.StringVar(&port, "port", "8081", "listening port")
}

func main() {
	flag.Parse()
	beego.RunMode = env

	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}

	models.ConnectDB()
	routers.Config()
	controllers.AppStart()
	beego.Run(":" + port)
}
