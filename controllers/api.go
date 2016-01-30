package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/gymer/pusher-api/models"
)

type httpError struct {
	Error string `json:"error"`
}

type APIController struct {
	beego.Controller
}

// @Title Push event
// @Description Push event with data to app specific channel
// @Param body    body  models.Event true    "body for event content"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /:appId/events [post]
func (u *APIController) Post() {
	var event models.Event
	appId := u.Ctx.Input.Params[":appId"]
	app := store.Apps[appId]

	if app == nil {
		httpResponseError(u.Ctx.ResponseWriter, 404, "Not found App with id "+appId)
		return
	}

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &event)
	if err != nil {
		httpResponseError(u.Ctx.ResponseWriter, 400, "Invalid JSON data")
		return
	}

	pushedClient := broadcastEvent(app, event)
	httpResponse(u.Ctx.ResponseWriter, 200, "{\"pushed_clients\": "+strconv.Itoa(pushedClient)+"}")
}

func init() {
	beego.InsertFilter("/v1/apps/*", beego.BeforeExec, apiAuthFilter)
}

func apiAuthFilter(ctx *context.Context) {
	if !apiAuthValidate(ctx) {
		ctx.ResponseWriter.Header().Set("WWW-Authenticate", `Basic realm="API realm"`)
		httpResponseError(ctx.ResponseWriter, 401, "Unauthorized")
	}
}

func apiAuthValidate(ctx *context.Context) bool {
	var app models.App
	r := ctx.Request
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	appId := ctx.Input.Params[":appId"]

	if len(s) != 2 || s[0] != "Basic" {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	err = models.DB.Where("id = ? and client_access_token = ? and server_access_token = ?", appId, pair[0], pair[1]).First(&app).Error
	if err != nil {
		return false
	}

	return true
}

func httpResponseError(w http.ResponseWriter, status int, message string) {
	error := httpError{message}
	b, _ := json.Marshal(error)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
	w.Write([]byte("\n"))
}

func httpResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(body))
	w.Write([]byte("\n"))
}
