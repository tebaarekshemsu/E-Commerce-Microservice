package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker hit",
	}
	
	_ = app.writeJSON(w, http.StatusOk, payload)
}