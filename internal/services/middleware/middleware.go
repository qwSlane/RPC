package middleware

import (
	"net/http"
)

type Middleware interface {
	IsAuthenticated(r *http.Request) bool
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleRegister(w http.ResponseWriter, r *http.Request)
}
