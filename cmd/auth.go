package main

import (
	"encoding/json"
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gobuffalo/envy"
)

type AuthData struct {
	AccessToken string `json="accessToken"`
	IDToken     string `json="id_token"`
}

type UserInfo struct {
	Login string `json="login"`
}

func main() {
	port := ":" + envy.Get("PORT", "3000")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		fakeToken := AuthData{AccessToken: "xxxxxxxyyyyyyyyyzzzz", IDToken: "rrrrrrrrrrrrrrrtttttttttttttttt"}
		b, _ := json.Marshal(fakeToken)
		w.Write(b)
	})

	r.Get("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		fakeUser := UserInfo{Login: "20009473"}
		b, _ := json.Marshal(fakeUser)
		w.Write(b)
	})

	log.Infof("======== App running in %v mode ========", "stage")
	http.ListenAndServe(port, r)
}
