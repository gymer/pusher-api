package main

import (
	"github.com/gymer/pusher-api/controllers"
	_ "github.com/gymer/pusher-api/docs"
	"github.com/gymer/pusher-api/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}

	routers.Config()
	controllers.AppStart()
	beego.Run()
}
