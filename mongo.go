package main

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

func MongoDB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("helpdesk"))
		defer s.Close()
		c.Next()
	}
}

const (
	UsersC          string = "users"
	SessionsC       string = "sessions"
	TicketsC        string = "tickets"
	TicketStatusC   string = "status"
	DomainsC        string = "domains"
	DocumentC       string = "documents"
	DomainSettingsC string = "domain_settings"
)
