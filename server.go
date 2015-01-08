package main

import (
	"github.com/go-martini/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
)

var CFG Config

func main() {
	if !CFG.Load() {
		if CFG.FirstTimeSetup() {
			CFG.Save()
		} else {
			return
		}
	} else {

	}

	m := martini.Classic()

	m.Use(martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/o"}))
	m.Use(MongoDB())

	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	InitTicketService(m)
	InitializeUserService(m)
	InitializeDepartmentService(m)
	InitTicketService(m)

	//m.Use(gzip.All())
	m.RunOnAddr(":80")

}
