package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"time"
)

func InitTicketService(m *martini.ClassicMartini) {
	m.Group("/o/ticket", func(r martini.Router) {
		r.Get("/info/:id", func(db *mgo.Database, params martini.Params) string {
			return params["id"]
		})
		r.Group("/list", func(r martini.Router) {
			r.Get("/total_count", RequireLogin(), func(db *mgo.Database) string {
				c := db.C(TicketsC)
				count, _ := c.Find(bson.M{}).Count()
				return strconv.Itoa(count)
			})
			r.Get("/mine", RequireLogin(), func(u User, db *mgo.Database) string {
				return HandleUserTickets(db, u, "")
			})
			r.Get("/:area/**", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {

				switch p["area"] {
				case "department":
					return HandleDepartmentTickets(db, u, p["_1"])
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

func HandleDepartmentTickets(db *mgo.Database, u User, t string) string {
	allow := false
	println(t)
	for i := 0; i < len(u.Department); i++ {
		if u.Department[i].Hex() == t {
			allow = true
			break
		}
	}

	if u.Roles.DomainAdmin {
		allow = true
	}

	if allow {
		c := db.C(TicketsC)
		var tickets []Ticket
		c.Find(bson.M{"department": bson.ObjectIdHex(t)}).All(&tickets)

		b, _ := json.Marshal(&tickets)
		return string(b)
	}
	return "error"
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
