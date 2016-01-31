// @APIVersion 1.0.0
// @Title Gymmer backend API
// @Description Gymmer backend API
// @Contact akrasman@gmail.com
// @TermsOfServiceUrl http://gymer.ws/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/gymer/pusher-api/controllers"

	"github.com/astaxie/beego"
)

func Config() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebsocketController{},
			),
		),
		beego.NSNamespace("/apps",
			beego.NSInclude(
				&controllers.EventsController{},
				&controllers.ChannelsController{},
			),
		),
	)
	beego.AddNamespace(ns)
	controllers.ApiFilters()
}
