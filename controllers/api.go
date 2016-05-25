package controllers

import (
	"encoding/json"
	"net/http"
)

type flatJson map[string]interface{}

type httpError struct {
	Error string `json:"error"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
}

func httpResponseJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(json)
}

func HttpResponseError(w http.ResponseWriter, status int, message string) {
	error := httpError{message}

	httpResponseJson(w, status, error)
}
