package main

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

func InitTicketService(m *martini.ClassicMartini) {
	m.Group("/o/ticket", func(r martini.Router) {
		r.Get("/:id", func(db *mgo.Database, params martini.Params) string {
			return params["id"]
		})
		r.Group("/list", func(r martini.Router) {
			r.Get("/", func() string {
				return ""
			})
			r.Get("/:area/:type", func(db *mgo.Database, p martini.Params) string {
				return p["area"] + "/" + p["type"]
			})
		})
	})
}
