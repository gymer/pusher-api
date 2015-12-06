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
curl --data '{"event": "new_message", "channel":"notifications", "data": {"title": "Hello", "content": "World"}}' --user CLIENT_ACCESS_TOKEN:SERVER_ACCESS_TOKEN localhost:8080/v1/app/APP_ID/events
```