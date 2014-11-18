package main

import (
	"github.com/go-martini/martini"
	//goauth2 "github.com/golang/oauth2"
	//"github.com/martini-contrib/gzip"
	//"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	//"labix.org/v2/mgo"
)

//var Cfg *oauth2.Config
//var MongoSession *mgo.Session

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
