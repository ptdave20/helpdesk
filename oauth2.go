package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/golang/oauth2"
	"github.com/golang/oauth2/google"
	"github.com/martini-contrib/sessions"
	"net/http"
)

func Oauth2Handler() martini.Handler {
	cfg, _ := google.NewConfig(&oauth2.Options{
		ClientID:     "812975936151-i8po4eflb6fggohokgl98d5998uh4t6k.apps.googleusercontent.com",
		ClientSecret: "5oqoK2q-_lnHO5kCdB8DjSyh",
		RedirectURL:  "http://localhost/o/token",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	})

	if err != nil {
		panic(err)
	}

	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/o/login":
				http.Redirect(w, r, cfg.AuthCodeURL("token", "offline", "force"), http.StatusOK)
				break
			case "/o/token":
				u := r.URL
				next := u.Query().Get("state")
				code := u.Query().Get("code")
				transport, err := cfg.NewTransportWithCode(code)
				if err != nil {
					http.Redirect(w, r, "/o/error", 302)
				}
				val, _ := json.Marshal(t.Token())
				s.Set("helpdesk", val)
				http.Redirect(w, r, next, 302)
				break
			case "/o/logout":
				u := r.URL
				next := u.Query().Get("next")
				s.Delete("helpdesk")
				http.Redirect(w, r, next, 302)
				break
			}
		}
	}

	c.Get("/o/login", func(w http.ResponseWriter, r *http.Request) {

	})

	c.Get("/o/token", func(w http.ResponseWriter, r *http.Request, s sessions.CookieStore) {
		u := r.URL
		code := u.Query().Get("code")
		if code != "" {
			transport, err := cfg.NewTransportWithCode(code)
			if err != nil {
				panic(err)
			}

		} else {
			return http.StatusBadRequest
		}
	})
}
