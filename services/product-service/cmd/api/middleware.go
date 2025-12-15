package main

import (
	"errors"
	"net/http"
	"strings"
)

func (app *Config) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.errorJSON(w, errors.New("missing authorization header"), http.StatusUnauthorized)
			return
		}

		// Split the header
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("invalid authorization header"), http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		// For demonstration, we consider "admin-token" as a valid admin token
		// In next Phase we will communicate with auth-service to validate tokens
		if token != "admin-token" {
			app.errorJSON(w, errors.New("unauthorized - admin access required"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
