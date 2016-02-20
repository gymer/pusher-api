package controllers

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/gymer/pusher-api/models"
)

type ApiBaseController struct {
	beego.Controller
}

type flatJson map[string]interface{}

type httpError struct {
	Error string `json:"error"`
}

func ApiFilters() {
	beego.InsertFilter("/v1/apps/*", beego.BeforeExec, apiAuthFilter)
}

func apiAuthFilter(ctx *context.Context) {
	if !apiAuthValidate(ctx) {
		ctx.ResponseWriter.Header().Set("WWW-Authenticate", `Basic realm="API realm"`)
		httpResponseError(ctx, 401, "Unauthorized")
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

	findOrAddApp(appId)

	return true
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
}

func httpResponseJson(ctx *context.Context, status int, data interface{}) {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(status)

	ctx.Output.Json(data, false, false)
}

func httpResponseError(ctx *context.Context, status int, message string) {
	error := httpError{message}

	httpResponseJson(ctx, status, error)
}

func (c *ApiBaseController) HttpResponseJson(status int, data interface{}) {
	httpResponseJson(c.Ctx, status, data)
}

func (c *ApiBaseController) HttpResponseError(status int, message string) {
	httpResponseError(c.Ctx, status, message)
}

func (c *ApiBaseController) app() *models.App {
	appId := c.Ctx.Input.Params[":appId"]
	return getApp(appId)
}
