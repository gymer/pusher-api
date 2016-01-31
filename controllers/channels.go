package controllers

type ChannelsController struct {
	ApiBaseController
}

// @Title GET channel
// @Description Fetch info for channel
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

// @Title GET channel users
// @Description Fetch all users subscribed to the channel
// @Success 200 body is empty
// @Failure 404 body is empty
// @router /:appId/channels/:channelName/users [get]
func (c *ChannelsController) GetUsers() {
	app := c.app()
	channelName := c.Ctx.Input.Params[":channelName"]
	subscrubers := app.ChannelSubscribers(channelName)
	usersSlice := make([]flatJson, len(subscrubers))
	i := 0

	for client, subscribed := range subscrubers {
		if subscribed {
			usersSlice[i] = flatJson{"id": client.Uuid}
			i++
		}
	}

	users := flatJson{"users": usersSlice}
	c.HttpResponseJson(200, users)
}
