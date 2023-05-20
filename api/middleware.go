package api

import (
	"net/http"
)

type Middleware interface {
	isAuthenticated(r *http.Request) bool
	handleLogin(w http.ResponseWriter, r *http.Request)
	handleRegister(w http.ResponseWriter, r *http.Request)
}
