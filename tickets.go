package main

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

func InitTicketService(m *martini.ClassicMartini) {
	m.Group("/o/ticket", func(r martini.Router) {
		r.Get("/info/:id", func(db *mgo.Database, params martini.Params) string {
			return params["id"]
		})
		r.Group("/list", func(r martini.Router) {
			r.Get("/**", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
				switch p["_1"] {
				case "department":
					return HandleDepartmentTickets(db, p["_2"])
					break
				case "mine":
					return HandleUserTickets(db, u, p["_2"])
					break
				case "user":
					return "unimplemented"
					break
				}
				return ""
			})
		})
		r.Post("/update/:id", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			return ""
		})
		r.Post("/insert", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			return ""
		})
	})
}

func HandleDepartmentTickets(db *mgo.Database, t string) string {
	return t
}

func HandleAssignedTickets(db *mgo.Database, t string) string {
	return t
}

func HandleUserTickets(db *mgo.Database, u User, t string) string {

	return t
}
