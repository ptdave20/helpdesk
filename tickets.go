package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	//"strconv"
	"time"
)

func InitTicketService(m *martini.ClassicMartini) {
	m.Group("/o/ticket", func(r martini.Router) {
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
		r.Get("/:id", RequireLogin(), func(u User, db *mgo.Database, req *http.Request, p martini.Params) string {
			var ticket Ticket
			var denied SimpleResult
			denied.Result = false
			c := db.C(TicketsC)
			err := c.Find(bson.M{"_id": p["id"]}).One(&ticket)
			if err != nil {
				return denied.Marshal()
			}

			// Can the user on it's own view the ticket ( based on submitter, assigned to, building, or department )?
			if u.CanViewTicket(ticket) {
				return ticket.Marshal()
			}

			// They did not have permissions or roles
			return denied.Marshal()
		})
		// If we post to a ticket id, we are trying to update
		r.Post("/:id", RequireLogin(), func(u User, db *mgo.Database, req *http.Request) string {
			//d := json.NewDecoder(req.Body)
			//var ticket Ticket
			return ""
		})
		// If a ticket is sent to us on the root of /o/ticket using POST, they are trying to add a ticket
		r.Post("/", RequireLogin(), func(u User, db *mgo.Database, req *http.Request) string {
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
	// This is strictly for retrieving tickets from the database
	m.Group("/o/tickets", func(r martini.Router) {
		r.Get("/", RequireLogin(), func() string {
			return "area required"
		})
		r.Get("/assigned", RequireLogin(), func() string {
			return "status required"
		})
		r.Get("/department", RequireLogin(), func() string {
			return "department required"
		})
		r.Get("/submitted", RequireLogin(), func() string {
			return "status required"
		})
		r.Get("/assigned/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var tickets []Ticket

			switch p["status"] {
			case "all":
				c.Find(bson.M{
					"assigned_to": u.Id}).All(&tickets)
				break
			case "open":
				c.Find(bson.M{
					"assigned_to": u.Id,
					"status":      bson.M{"$ne": "closed"}}).All(&tickets)
				break
			default:
				c.Find(bson.M{
					"assigned_to": u.Id,
					"status":      p["status"]}).All(&tickets)
				break
			}

			c.Find(bson.M{"assigned_to": u.Id}).All(&tickets)

			b, _ := json.Marshal(tickets)

			return string(b)
		})
		r.Get("/department/:id/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var tickets []Ticket
			if !u.InDepartment(bson.ObjectId(p["id"])) {
				// We do nothing
			} else {
				switch p["status"] {
				case "all":
					c.Find(bson.M{
						"department": bson.ObjectId(p["id"])}).All(&tickets)
					break
				case "open":
					c.Find(bson.M{
						"department": bson.ObjectId(p["id"]),
						"status":     bson.M{"$ne": "closed"}}).All(&tickets)
					break
				default:
					c.Find(bson.M{
						"department": bson.ObjectId(p["id"]),
						"status":     p["status"]}).All(&tickets)
					break
				}

			}

			b, _ := json.Marshal(tickets)

			return string(b)
		})
		r.Get("/submitted/:status", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var tickets []Ticket

			switch p["status"] {
			case "all":
				c.Find(bson.M{
					"submitter": u.Id}).All(&tickets)
				break
			case "open":
				c.Find(bson.M{
					"submitter": u.Id,
					"status":    bson.M{"$ne": "closed"}}).All(&tickets)
				break
			default:
				c.Find(bson.M{
					"submitter": u.Id,
					"status":    p["status"]}).All(&tickets)
				break
			}
			b, _ := json.Marshal(tickets)

			return string(b)
		})
	})
}
