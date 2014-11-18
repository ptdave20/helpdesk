package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()

	m.Use(MongoDB())
	//m.Use(gzip.All())
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	m.Get("/data", RequireLogin(), func(user *User) string {
		return "Hello " + user.Firstname
	})
	InitTicketService(m)
	m.RunOnAddr(":80")

}
