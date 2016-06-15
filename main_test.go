package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/gymer/pusher-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	localAddress             string = "localhost:3000"
	invalidClientAccessToken string = "wrong-access-token"
	validClientAccessToken   string = "7f075634bdc4d2ef"
	validServerAccessToken   string = "97c1aa41fc71b92e"
	invalidAppID             int    = 124141
	validAppID               int    = 1
)

func wsUrl(path string) string {
	return "ws://" + localAddress + "/v1/ws" + path
}

func httpUrl(path string) string {
	return "http://" + localAddress + "/v1" + path
}

// TestGet is a sample to run an endpoint test
func TestWS(t *testing.T) {

	Convey("Subject: Wrong app code\n", t, func() {
		url := wsUrl("/app/" + invalidClientAccessToken)
		ws, r, err := websocket.DefaultDialer.Dial(url, nil)

		_, _, err = ws.ReadMessage()
		defer ws.Close()

		Convey("Status Code Should Be 101", func() {
			So(r.StatusCode, ShouldEqual, 101)
		})

		Convey("Recive WS close message", func() {
			So(err.Error(), ShouldEqual, "websocket: close 4001 Invalid app code")
		})
	})

	Convey("Subject: Valid app code\n", t, func() {
		var event models.Event

		url := wsUrl("/app/" + validClientAccessToken)
		ws, r, err := websocket.DefaultDialer.Dial(url, nil)
		defer ws.Close()

		_, b, err := ws.ReadMessage()
		json.Unmarshal(b, &event)

		Convey("Status Code Should Be 101", func() {
			So(r.StatusCode, ShouldEqual, 101)
		})

		Convey("Recive WS close message", func() {
			So(err, ShouldBeNil)
			So(event.Name, ShouldEqual, "gymer:connection_established")
		})

		Convey("Subscribe to public channel", func() {
			event := models.Event{Name: "gymer:subscribe", Channel: "notifications"}
			err := ws.WriteJSON(event)
			httpClient := &http.Client{}

			_, b, err = ws.ReadMessage()

			So(err, ShouldBeNil)

			Convey("Recive success subscribtion message", func() {
				var event models.Event
				json.Unmarshal(b, &event)

				So(event.Name, ShouldEqual, "gymer:subscription_success")
			})

			Convey("Recive push event via API request", func() {
				var jsonStr = []byte(`{"event": "new_message", "channel":"notifications", "data": {"title": "Hello", "content": "World"}}`)

				req, _ := http.NewRequest("POST", httpUrl("/apps/"+strconv.Itoa(validAppID)+"/events"), bytes.NewBuffer(jsonStr))
				req.Header.Set("Content-Type", "application/json")

				Convey("With invalid basic auth", func() {
					res, _ := httpClient.Do(req)
					body, _ := ioutil.ReadAll(res.Body)

					Convey("Should block request", func() {
						So(res.StatusCode, ShouldEqual, 401)
						So(string(body), ShouldEqual, `{"error":"Unauthorized"}`)
						res.Body.Close()
					})

				})

				Convey("With valid basic auth", func() {
					req.SetBasicAuth(validClientAccessToken, validServerAccessToken)
					res, _ := httpClient.Do(req)
					body, _ := ioutil.ReadAll(res.Body)

					if res.StatusCode != http.StatusOK {
						t.Fatalf("Push request failed: %+v", res)
					}

					Convey("Should Pass request", func() {
						So(res.StatusCode, ShouldEqual, 200)
						So(string(body), ShouldEqual, `{"pushed_clients":1}`)
						res.Body.Close()
					})

					Convey("WS client recieve push", func() {
						var event models.Event
						_, b, err := ws.ReadMessage()

						json.Unmarshal(b, &event)

						So(err, ShouldBeNil)
						So(event.Name, ShouldEqual, "new_message")
						So(event.Data["title"], ShouldEqual, "Hello")
						So(event.Data["content"], ShouldEqual, "World")
					})

				})

			})
		})
	})
}
