package main

import (
	"github.com/go-martini/martini"
	//goauth2 "github.com/golang/oauth2"
	//"github.com/martini-contrib/gzip"
<<<<<<< HEAD
	"github.com/martini-contrib/oauth2"
=======
	//"github.com/martini-contrib/oauth2"
>>>>>>> 52179a7d6a38c1c9f904c475275d43bf740a9667
	"github.com/martini-contrib/sessions"
	//"labix.org/v2/mgo"
)

//var Cfg *oauth2.Config
//var MongoSession *mgo.Session

func main() {
	m := martini.Classic()
<<<<<<< HEAD
	oauth2.PathLogin = "/o/login"
	oauth2.PathLogout = "/o/logout"
	//oauth2.PathCallback = "/o/token"
	oauth2.PathError = "/o/error"

	//m.Use(MongoDB())
	//m.Use(gzip.All())
	m.Use(sessions.Sessions("helpdesk-session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())
	/*
			m.Use(oauth2.Google(
				goauth2.Client("812975936151-i8po4eflb6fggohokgl98d5998uh4t6k.apps.googleusercontent.com", "5oqoK2q-_lnHO5kCdB8DjSyh"),
				goauth2.RedirectURL("http://localhost/o/token"),
				goauth2.Scope("https://www.googleapis.com/auth/userinfo.profile"),
			))

		// Tokens are injected to the handlers
		m.Get("/", func(tokens oauth2.Tokens) string {
			if tokens.Expired() {
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
	*/
=======

	m.Use(MongoDB())
	//m.Use(gzip.All())
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	m.Get("/data", RequireLogin(), func(user *User) string {
		return "Hello " + user.Firstname
	})
<<<<<<< HEAD

>>>>>>> 52179a7d6a38c1c9f904c475275d43bf740a9667
=======
	InitTicketService(m)
>>>>>>> 0116bb3a09ddc21a241e3a4f1df67a3bd1cd19f3
	m.RunOnAddr(":80")

}
