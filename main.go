package main

import (
	_ "github.com/gymer/pusher-api/docs"
	_ "github.com/gymer/pusher-api/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
