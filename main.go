package main

import (
	"github.com/gymer/pusher-api/controllers"
	_ "github.com/gymer/pusher-api/docs"
	"github.com/gymer/pusher-api/models"
	"github.com/gymer/pusher-api/routers"

	"github.com/astaxie/beego"
	"github.com/namsral/flag"
)

const defaultPort = "3001"

var (
	env, port string
	err       error
)

func init() {
	flag.StringVar(&env, "env", "dev", "set app environment")
	flag.StringVar(&port, "port", defaultPort, "listening port")
}

func main() {
	flag.Parse()
	beego.RunMode = env

	if beego.RunMode == "dev" {
		beego.EnableAdmin = true
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}

	startServer()
}

func startServer() {
	models.ConnectDB()
	routers.Config()
	controllers.AppStart()
	beego.Run(":" + port)
}
