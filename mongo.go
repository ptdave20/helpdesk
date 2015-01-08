package main

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
)

func MongoDB() martini.Handler {
	session, err := mgo.Dial(CFG.MongoAddress)
	if err != nil {
		println(CFG.MongoAddress)
		log.Fatal(err)
	}
	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB(CFG.MongoDatabase))
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
