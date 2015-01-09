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
	m.Get("/o/my_departments", RequireLogin(), func(u User, d Domain, db *mgo.Database, p martini.Params) string {
		var deps []Department
		for i := 0; i < len(d.Departments); i++ {
			if d.Departments[i].IsBuildingSpecific {
				if (d.Departments[i].Building.Hex() == u.Building.Hex()) || u.Roles.DomainAdmin {
					deps = append(deps, d.Departments[i])
				}
			} else {
				deps = append(deps, d.Departments[i])
			}
		}
		b, _ := json.Marshal(deps)
		return string(b)
	})
	m.Group("/o/ticket", func(r martini.Router) {
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
			bId := bson.ObjectIdHex(p["id"])
			if !bId.Valid() {
				return "invalid"
			}
			err := c.Find(bson.M{"_id": bson.ObjectIdHex(p["id"])}).One(&ticket)
			if err != nil {
				print(err.Error())
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
		r.Post("/:id", RequireLogin(), func(u User, db *mgo.Database, req *http.Request, p martini.Params) string {
			c := db.C(TicketsC)
			var original Ticket
			id := bson.ObjectIdHex(p["id"])
			if !id.Valid() {
				return "invalid"
			}

			c.Find(bson.M{"_id": id}).One(&original)

			if !u.CanEditTicket(original) {
				return "denied"
			}

			d := json.NewDecoder(req.Body)
			var ticket TicketUpdate
			err := d.Decode(&ticket)
			if err != nil {
				panic(err)
			}

			var changes = false

			// Editing assumes description, subject, department and category
			if u.CanEditTicket(original) {
				original.Department = ticket.Department
				original.Category = ticket.Category
				original.Subject = ticket.Subject
				original.Description = ticket.Description
				original.Domain = u.Domain
				changes = true
			}

			if changes {
				// Set a new update time
				original.Updated = time.Now()

				err = c.Update(bson.M{"_id": id}, original)
				if err != nil {
					panic(err)

				}
				return "success"
			}
			return "no changes made"
		})
		r.Post("/:id/closed", RequireLogin(), func(u User, db *mgo.Database, req *http.Request, p martini.Params) string {
			return "unimplemented"
		})
		// If a ticket is sent to us on the root of /o/ticket using POST, they are trying to add a ticket
		r.Post("", RequireLogin(), func(domain Domain, u User, db *mgo.Database, req *http.Request) string {
			d := json.NewDecoder(req.Body)

			var tkt Ticket

			err := d.Decode(&tkt)
			if err != nil {
				panic(err)
			}

			tkt.Id = bson.NewObjectId()
			tkt.Submitter = u.Id
			tkt.Status = "open"
			tkt.Building = u.Building
			tkt.Domain = domain.Id
			tkt.Created = time.Now()

			c := db.C(TicketsC)
			err = c.Insert(tkt)
			if err != nil {
				panic(err)
			}
			return "success"
		})
	})
	// This is strictly for retrieving tickets from the database
	m.Group("/o/tickets", func(r martini.Router) {
		r.Get("/", RequireLogin(), func() string {
			return "area required"
		})
		r.Get("/building", RequireLogin(), func(u User, db *mgo.Database) string {
			if (u.Roles.BldgViewTicket || u.Roles.DomainAdmin) && u.Building.Hex() != "" {
				c := db.C(TicketsC)
				var tickets []Ticket
				c.Find(bson.M{"building": u.Building}).All(&tickets)
				b, _ := json.Marshal(tickets)
				return string(b)
			} else {
				if u.Building.Hex() == "" {
					return "no building set for user"
				} else if !u.Roles.BldgViewTicket {
					return "not allowed to see tickets for building"
				}
			}
			return "unknown error"
		})
		r.Get("/building", RequireLogin(), func(u User, db *mgo.Database) string {
			if !u.Building.Valid() {
				return "invalid building"
			}

			c := db.C(TicketsC)

			var tickets []Ticket

			c.Find(bson.M{"building": u.Building}).All(&tickets)
			b, _ := json.Marshal(tickets)
			return string(b)
		})
		r.Get("/building/:id", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			if !u.Roles.DomainAdmin || (u.Building.Hex() != p["id"] && u.Roles.BldgViewTicket) {
				return "access denied"
			}
			c := db.C(TicketsC)
			var tickets []Ticket
			c.Find(bson.M{"building": bson.ObjectIdHex(p["id"])}).All(&tickets)
			b, _ := json.Marshal(tickets)
			return string(b)
		})
		r.Get("/assigned", RequireLogin(), func(u User, db *mgo.Database) string {
			c := db.C(TicketsC)
			count, _ := c.Find(bson.M{"assigned_to": u.Id, "status": bson.M{"$ne": "closed"}}).Count()
			return strconv.Itoa(count)
		})
		r.Get("/department", RequireLogin(), func() string {
			return "department required"
		})
		r.Get("/submitted", RequireLogin(), func(u User, db *mgo.Database) string {
			c := db.C(TicketsC)
			count, _ := c.Find(bson.M{"submitter": u.Id, "status": bson.M{"$ne": "closed"}}).Count()
			return strconv.Itoa(count)
		})
		r.Get("/assigned/:status/:count", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var tickets []Ticket
			var count int
			count, _ = strconv.Atoi(p["count"])
			switch p["status"] {
			case "all":
				c.Find(bson.M{
					"assigned_to": u.Id}).Limit(count).All(&tickets)
				break
			case "open":
				c.Find(bson.M{
					"assigned_to": u.Id,
					"status":      bson.M{"$ne": "closed"}}).Limit(count).All(&tickets)
				break
			default:
				c.Find(bson.M{
					"assigned_to": u.Id,
					"status":      p["status"]}).Limit(count).All(&tickets)
				break
			}

			c.Find(bson.M{"assigned_to": u.Id}).All(&tickets)

			b, _ := json.Marshal(tickets)

			return string(b)
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
		r.Get("/department/:id", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var id bson.ObjectId = bson.ObjectIdHex(p["id"])
			if !id.Valid() {
				return "invalid department"
			}
			if !u.InDepartment(id) || !u.Roles.DomainAdmin {
				return "denied"
			}

			count, _ := c.Find(bson.M{"department": id, "status": bson.M{"$ne": "closed"}}).Count()
			return strconv.Itoa(count)
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
		r.Get("/submitted/:status/:count", RequireLogin(), func(u User, db *mgo.Database, p martini.Params) string {
			c := db.C(TicketsC)
			var tickets []Ticket
			var count int
			count, _ = strconv.Atoi(p["count"])
			switch p["status"] {
			case "all":
				c.Find(bson.M{
					"submitter": u.Id}).Limit(count).All(&tickets)
				break
			case "open":
				c.Find(bson.M{
					"submitter": u.Id,
					"status":    bson.M{"$ne": "closed"}}).Limit(count).All(&tickets)
				break
			default:
				c.Find(bson.M{
					"submitter": u.Id,
					"status":    p["status"]}).Limit(count).All(&tickets)
				break
			}

			b, _ := json.Marshal(tickets)

			return string(b)
		})
	})
}
