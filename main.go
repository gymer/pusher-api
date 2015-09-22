package main

import (
	_ "github.com/Iverson/pusher-api/docs"
	_ "github.com/Iverson/pusher-api/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
