# pusher-api

## Run

```
bee run
```
In browser connect to
```
ws://localhost:8080/v1/ws/app/secret-app-key`
```
In terminal
```
curl --data '{"event": "new_message", "channel":"notifications", "data": {"title": "Hello", "content": "World"}}' --user 8cb981cfacddfe3c:SERVER_ACCESS_TOKEN localhost:8080/v1/apps/APP_ID/events
```

## Build

Ubuntu 14.04
```
GOOS=linux GOARCH=amd64 go build
```