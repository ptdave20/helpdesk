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
<<<<<<< HEAD
=======

const (
	UsersC          string = "users"
	SessionsC       string = "sessions"
	TicketsC        string = "tickets"
	DomainsC        string = "domains"
	DepartmentsC    string = "departments"
	DocumentC       string = "documents"
	DomainSettingsC string = "domain_settings"
)
>>>>>>> 52179a7d6a38c1c9f904c475275d43bf740a9667
