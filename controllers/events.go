package controllers

import (
	"encoding/json"

	"github.com/gymer/pusher-api/models"
)

type EventsController struct {
	ApiBaseController
}

// @Title Push event
// @Description Push event with data to app specific channel
// @Param body    body  models.Event true    "body for event content"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /:appId/events [post]
func (c *EventsController) Post() {
	var event models.Event

	app := c.app()

	if app == nil {
		c.HttpResponseError(404, "Not found App")
		return
	}

	Logger.Warn("POST bytes: %+v \n", c.Ctx.Input.RequestBody)
	Logger.Warn("POST string: %+v \n", string(c.Ctx.Input.RequestBody[:]))

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &event)
	if err != nil {
		c.HttpResponseError(400, "Invalid JSON data")
		return
	}

	pushedClient := broadcastEvent(app, event)
	resp := make(map[string]interface{})
	resp["pushed_clients"] = pushedClient

	c.HttpResponseJson(200, resp)
}
