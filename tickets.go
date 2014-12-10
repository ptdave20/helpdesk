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
			r.Get("/mine/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
				return HandleUserTickets(db, u, p["status"])
			})
			r.Get("/:area/:id/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
				switch p["area"] {
				case "department":
					return HandleDepartmentTickets(db, u, p["id"], p["status"])
					break
				case "user":
					return "unimplemented"
					break
				}
				return ""
			})
		})
		r.Group("/count", func(r martini.Router) {
			r.Get("/department/:id/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
				found := false
				for i := 0; i < len(u.Department); i++ {
					if p["id"] == u.Department[i].Hex() {
						found = true
						break
					}
				}
				if !found && !u.Roles.DomainAdmin {
					return "denied"
				}

				c := db.C(TicketsC)
				if p["status"] == "open" {
					count, err := c.Find(bson.M{"department": bson.ObjectIdHex(p["id"]), "status": bson.M{"$ne": "closed"}}).Count()
					if err != nil {
						panic(err)
					}

					return strconv.Itoa(count)
				}
				count, err := c.Find(bson.M{"department": bson.ObjectIdHex(p["id"]), "status": p["status"]}).Count()
				if err != nil {
					panic(err)
				}

				return strconv.Itoa(count)
			})
		})

		r.Post("/update/:id", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			return ""
		})
		r.Post("/close/:id", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			var id string = p["id"]
			var tkt Ticket

			c := db.C(TicketsC)
			err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&tkt)
			if err != nil {
				panic(err)
			}

			var allow bool = false

			// Check to see if the user is allowed to close this based on ownership
			if u.Id.Hex() == tkt.Submitter.Hex() {
				allow = true
			}

			// Check to see if the user is allowed to close this based on assignment
			if u.Id.Hex() == tkt.AssignedTo.Hex() {
				allow = true
			}

			// Check to see if the user is allowed to close this based on
			if u.Roles.DomainAdmin {
				allow = true
			}

			var res SimpleResult
			if allow {
				err := c.Update(bson.M{"_id": tkt.Id}, bson.M{"closed": time.Now(), "status": "closed"})
				if err != nil {
					panic(err)
				}
				res.Result = true
			} else {
				res.Result = false
			}

			b, _ := json.Marshal(res)
			return string(b)
		})
		// If we post to a ticket id, we are trying to update
		r.Post("/:id", RequireLogin(), func(u User, db *mgo.Database, req *http.Request) {
			d := json.NewDecoder(req.Body)
			var ticket Ticket
		})
		// If a ticket is sent to us on the root of /o/ticket using POST, they are trying to add a ticket
		r.Post("/", RequireLogin(), func(u User, db *mgo.Database, req *http.Request) {
			d := json.NewDecoder(req.Body)

			var tkt Ticket

			err := d.Decode(&tkt)
			if err != nil {
				panic(err)
			}

			tkt.Id = bson.NewObjectId()
			tkt.Submitter = u.Id
			tkt.Status = "open"

			tkt.Created = time.Now()

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

func HandleDepartmentTickets(db *mgo.Database, u User, t string, s string) string {
	allow := false
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
		if s != "closed" {
			c.Find(bson.M{"department": bson.ObjectIdHex(t), "status": bson.M{"$ne": "closed"}}).All(&tickets)
		} else {
			if s == "all" {
				c.Find(bson.M{"department": bson.ObjectIdHex(t)}).All(&tickets)
			} else {
				c.Find(bson.M{"department": bson.ObjectIdHex(t), "status": "closed"}).All(&tickets)
			}
		}
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

	if t != "closed" {
		c.Find(bson.M{"submitter": u.Id, "status": bson.M{"$ne": "closed"}}).All(&tkts)

	} else if t == "all" {
		c.Find(bson.M{"submitter": u.Id}).All(&tkts)
	} else if t == "closed" {
		c.Find(bson.M{"submitter": u.Id, "status": "closed"}).All(&tkts)
	}

	j, _ := json.Marshal(&tkts)

	return string(j)
}
