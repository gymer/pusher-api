package controllers

type ChannelsController struct {
	ApiBaseController
}

// @Title GET channel
// @Description Fetch info for channel
// @Param body    body  models.Event true    "body for event content"
// @Success 200 body is empty
// @Failure 404 body is empty
// @router /:appId/channels/:channelName [get]
func (c *ChannelsController) Get() {
	app := c.app()
	channelName := c.Ctx.Input.Params[":channelName"]
	subscrubers := app.ChannelSubscribers(channelName)
	channel := flatJson{"used": false, "subscribers": 0}

	if subscrubers != nil {
		channel["used"] = true
		channel["subscribers"] = len(subscrubers)
	}

	c.HttpResponseJson(200, channel)
}
