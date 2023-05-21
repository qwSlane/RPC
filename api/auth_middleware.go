package api

import (
	"encoding/json"
	"errors"
	"log"
	"main/storage"
	"main/types"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signinKey = "qpowuD63DS&hfht91"
)

type AuthMiddleware struct {
	Storage storage.Storage
}

func NewAuthMiddleware(storage storage.Storage) *AuthMiddleware {
	return &AuthMiddleware{
		Storage: storage,
	}
}

func (s *AuthMiddleware) isAuthenticated(r *http.Request) bool {

	tokenString := r.Header.Get("Authorization")

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signin method ")
		}
		return []byte(signinKey), nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if validUser := s.Storage.GetUser(username); validUser == nil || !ok {
			log.Println("Invalid username in token")
			return false
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			log.Println("Invalid expiration time in token")
			return false
		}

		if int64(exp) < time.Now().Unix() {
			log.Println("Token has expired")
			return false
		}

		return true
	}

	return false
}

func (s *AuthMiddleware) handleLogin(w http.ResponseWriter, r *http.Request) {

	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	validUser := s.Storage.CheckUser(user.Username, user.Password)

	tokenString, err := generateToken(*validUser)
	if err != nil {
		http.Error(w, "Couldn't generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", tokenString)
	if err != nil {
		http.Error(w, "Could not send token", http.StatusInternalServerError)
		return
	}

	log.Printf("User: %s was authorized", validUser.Username)

}

func (s *AuthMiddleware) handleRegister(w http.ResponseWriter, r *http.Request) {

	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	isOriginal := s.Storage.GetUser(user.Username)

	if isOriginal != nil {
		http.Error(w, "User already exist", http.StatusBadRequest)
		return
	}

	err = s.Storage.CreateUser(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))

}

func generateToken(user types.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(signinKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
