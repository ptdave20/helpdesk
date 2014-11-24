package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
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
		r.Post("/insert", RequireLogin(), func(u User, db *mgo.Database, req *http.Request) string {
			println("Trying to add ticket")

			d := json.NewDecoder(req.Body)

			var tkt Ticket

			err := d.Decode(&tkt)
			if err != nil {
				panic(err)
			}

			tkt.Id = bson.NewObjectId()
			tkt.Submitter = u.Id

			tkt.Created = time.Now()
			//tkt.Closed = nil

			c := db.C(TicketsC)
			err = c.Insert(tkt)
			if err != nil {
				panic(err)
			}

			// Move it to json
			out, err := json.Marshal(&tkt)

			if err != nil {
				panic(err)
			}
			return string(out)
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
	c := db.C(TicketsC)

	var tkts []Ticket

	c.Find(bson.M{"submitter": u.Id}).All(&tkts)

	j, _ := json.Marshal(&tkts)

	return string(j)
}
