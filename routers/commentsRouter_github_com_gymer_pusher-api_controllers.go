package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:APIController"] = append(beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:APIController"],
		beego.ControllerComments{
			"Post",
			`/:appId/events`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:WebsocketController"] = append(beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:WebsocketController"],
		beego.ControllerComments{
			"Connect",
			`/app/:key`,
			[]string{"get"},
			nil})

}
