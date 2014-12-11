package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
	"log"
)

var CFG Config

func main() {
	cfgBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(cfgBytes, &CFG)

	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()
	m.Use(martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/o"}))
	m.Use(MongoDB())
	m.Use(gzip.All())
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	InitTicketService(m)
	InitializeUserService(m)
	InitializeDepartmentService(m)
	InitTicketService(m)

	m.RunOnAddr(":80")

}
