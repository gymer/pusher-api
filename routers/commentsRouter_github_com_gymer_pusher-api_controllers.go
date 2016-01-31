package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:ChannelsController"] = append(beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:ChannelsController"],
		beego.ControllerComments{
			"Get",
			`/:appId/channels/:channelName`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:EventsController"] = append(beego.GlobalControllerRouter["github.com/gymer/pusher-api/controllers:EventsController"],
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
