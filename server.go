package main

import (
	"github.com/go-martini/martini"
	goauth2 "github.com/golang/oauth2"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"log"
)

//var Cfg *oauth2.Config
//var MongoSession *mgo.Session

func main() {
	//ConfigureOauth2()
	session, err := mgo.Dial("localhost")

	m := martini.Classic()
	oauth2.PathLogin = "/o/login"
	oauth2.PathLogout = "/o/logout"
	oauth2.PathCallback = "/o/token"
	oauth2.PathError = "/o/error"

	if err != nil {
		log.Fatal(err)
	} else {
		m.Map(session)
	}

	m.Use(gzip.All())
	m.Use(sessions.Sessions("helpdes-session", sessions.NewCookieStore([]byte("session"))))
	m.Use(oauth2.Google(&goauth2.Options{
		ClientID:     "812975936151-i8po4eflb6fggohokgl98d5998uh4t6k.apps.googleusercontent.com",
		ClientSecret: "5oqoK2q-_lnHO5kCdB8DjSyh",
		RedirectURL:  "http://localhost/o/token",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}))

	// Tokens are injected to the handlers
	m.Get("/", func(tokens oauth2.Tokens) string {
		if tokens.IsExpired() {
			return "not logged in, or the access token is expired"
		}
		return "logged in"
	})

	// Routes that require a logged in user
	// can be protected with oauth2.LoginRequired handler.
	// If the user is not authenticated, they will be
	// redirected to the login path.
	m.Get("/restrict", oauth2.LoginRequired, func(tokens oauth2.Tokens) string {
		return tokens.Access()
	})

	m.RunOnAddr(":80")
	//http.HandleFunc("/o/login", Oauth2Login)
	//http.HandleFunc("/o/token", Oauth2Token)

}
