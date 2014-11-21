package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/o"}))
	m.Use(MongoDB())
	m.Use(gzip.All())
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	m.Get("/data", RequireLogin(), func(user *User) string {
		return "Hello " + user.Firstname
	})
	InitTicketService(m)
	InitializeUserService(m)
	InitializeDepartmentService(m)

	m.RunOnAddr(":80")

}
