package main

import (
	"github.com/go-martini/martini"

	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func InitializeDepartmentService(m *martini.ClassicMartini) {
	m.Group("/o/department", func(r martini.Router) {
		r.Get("/list", RequireLogin(), func(domain Domain, db *mgo.Database) string {
			b, _ := json.Marshal(domain.Departments)
			return string(b)
		})
		r.Post("", RequireLogin(), func(domain Domain, u User, req *http.Request, db *mgo.Database) string {
			if !u.Roles.DomainAdmin || !u.Roles.DomainModDep {
				return "denied"
			}

			decoder := json.NewDecoder(req.Body)
			var d Department
			err := decoder.Decode(&d)
			if err != nil {
				panic(err)
			}

			d.Id = bson.NewObjectId()

			for i := 0; i < len(d.Category); i++ {
				d.Category[i].Id = bson.NewObjectId()
			}

			c := db.C(DomainsC)
			err = c.Update(bson.M{"_id": domain.Id}, bson.M{"$push": bson.M{"departments": d}})

			if err != nil {
				panic(err)
			}
			return "success"
		})
		r.Get("/:id", RequireLogin(), func(domain Domain, u User, req *http.Request, p martini.Params) string {
			for i := 0; i < len(domain.Departments); i++ {
				if p["id"] == domain.Departments[i].Id.Hex() {
					b, _ := json.Marshal(domain.Departments[i])
					return string(b)
				}
			}
			return "not found"
		})
		r.Post("/:id", RequireLogin(), func(domain Domain, u User, req *http.Request, db *mgo.Database, p martini.Params) string {
			if !u.Roles.DomainAdmin || !u.Roles.DomainModDep {
				return "denied"
			}
			id := bson.ObjectIdHex(p["id"])
			if !id.Valid() {
				return "invalid id"
			}

			decoder := json.NewDecoder(req.Body)
			var cat Category
			err := decoder.Decode(&cat)
			cat.Id = bson.NewObjectId()

			if err != nil {
				panic(err)
			}

			c := db.C(DomainsC)
			for i := 0; i < len(domain.Departments); i++ {
				if domain.Departments[i].Id.Hex() == id.Hex() {
					domain.Departments[i].Category = append(domain.Departments[i].Category, cat)
					break
				}
			}

			err = c.Update(bson.M{"_id": domain.Id}, domain)
			if err != nil {
				panic(err)
			}

			return "success"
		})
		r.Delete("/:id", RequireLogin(), func(domain Domain, u User, db *mgo.Database, p martini.Params) string {
			id := bson.ObjectIdHex(p["id"])
			if !id.Valid() {
				return "invalid id"
			}

			if !u.Roles.DomainAdmin || !u.Roles.DomainModDep {
				return "denied"
			}

			c := db.C(DomainsC)

			err := c.Update(bson.M{"_id": domain.Id}, bson.M{"$pull": bson.M{"departments": bson.M{"_id": id}}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
		r.Delete("/:id/:cat", RequireLogin(), func(domain Domain, u User, db *mgo.Database, p martini.Params) string {
			id := bson.ObjectIdHex(p["id"])
			cat := bson.ObjectIdHex(p["cat"])
			if !id.Valid() || !cat.Valid() {
				return "invalid id"
			}

			if !u.Roles.DomainAdmin || !u.Roles.DomainModDep {
				return "denied"
			}

			c := db.C(DomainsC)

			err := c.Update(bson.M{"_id": domain.Id}, bson.M{"$pull": bson.M{"departments.$.category": bson.M{"_id": cat}}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
	})
}
