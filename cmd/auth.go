package main

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/envy"
)

var hmacSampleSecret []byte

type AuthData struct {
	AccessToken string `json:"accessToken"`
	IDToken     string `json:"id_token"`
}

type UserInfo struct {
	Login string `json="login"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	password := r.URL.Query().Get("password")
	if password == "" {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "digital-data-platform",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
		"sub":      "digital-data-platform",
		"aud":      "digital-data-platform",
		"user":     "20009473",
		"username": "fclaeys",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("monsupersecretvraimentsupersecret"))
	if err != nil {
		log.WithError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func main() {

	port := ":" + envy.Get("PORT", "3000")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/login", LoginHandler)

	log.Infof("======== App running in %v mode ========", "stage")
	http.ListenAndServe(port, r)
}
