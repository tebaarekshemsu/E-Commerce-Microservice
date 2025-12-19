package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker hit",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) ProductServiceProxy() http.Handler {
	target, _ := url.Parse("http://product-service")
	return httputil.NewSingleHostReverseProxy(target)
}

func (app *Config) PaymentServiceProxy() http.Handler {
	target, _ := url.Parse("http://payment-service")
	return httputil.NewSingleHostReverseProxy(target)
}

func (app *Config) OrderServiceProxy() http.Handler {
	target, _ := url.Parse("http://order-service")
	return httputil.NewSingleHostReverseProxy(target)
}