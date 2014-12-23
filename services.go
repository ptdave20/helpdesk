package main

import (
	"github.com/go-martini/martini"

	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func InitializeDepartmentService(m *martini.ClassicMartini) {
	m.Group("/o/domain", func(r martini.Router) {
		r.Get("/buildings", RequireLogin(), func(domain Domain) string {
			b, _ := json.Marshal(domain.Buildings)
			return string(b)
		})
		r.Post("/building", RequireLogin(), func(domain Domain, user User, db *mgo.Database, req *http.Request) string {
			if !user.Roles.DomainAdmin || !user.Roles.DomainModBldg {
				return "denied"
			}
			decoder := json.NewDecoder(req.Body)
			var building Building
			decoder.Decode(&building)
			building.Id = bson.NewObjectId()
			c := db.C(DomainsC)
			err := c.UpdateId(domain.Id, bson.M{"$push": bson.M{"buildings": building}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
		r.Post("/building/:id", RequireLogin(), func(domain Domain, user User, db *mgo.Database, req *http.Request, p martini.Params) string {
			if !user.Roles.DomainAdmin || !user.Roles.DomainModBldg {
				return "denied"
			}
			id := bson.ObjectIdHex(p["id"])
			var bldg Building
			dec := json.NewDecoder(req.Body)
			dec.Decode(&bldg)
			c := db.C(DomainsC)
			err := c.Update(bson.M{"_id": domain.Id, "buildings._id": id}, bson.M{"$set": bson.M{"buildings.$": bldg}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
		r.Delete("/building/:id", RequireLogin(), func(domain Domain, user User, db *mgo.Database, req *http.Request, p martini.Params) string {
			if !user.Roles.DomainAdmin || !user.Roles.DomainModBldg {
				return "denied"
			}
			id := bson.ObjectIdHex(p["id"])
			c := db.C(DomainsC)
			err := c.UpdateId(domain.Id, bson.M{"$pull": bson.M{"buildings": bson.M{"_id": id}}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
	})
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

			// Pull in the incoming department data
			decoder := json.NewDecoder(req.Body)
			var incDep DepartmentUpdate
			decoder.Decode(&incDep)

			c := db.C(DomainsC)
			if incDep.Building.Valid() {
				err := c.Update(
					bson.M{"_id": domain.Id, "departments._id": id},
					bson.M{"$set": bson.M{"departments.$.name": incDep.Name, "departments.$.is_building_specific": incDep.IsBuildingSpecific, "departments.$.building": incDep.Building}})
				if err != nil {
					panic(err)
				}
				return "success"
			} else {
				err := c.Update(
					bson.M{"_id": domain.Id, "departments._id": id},
					bson.M{"$set": bson.M{"departments.$.name": incDep.Name, "departments.$.is_building_specific": incDep.IsBuildingSpecific}})
				if err != nil {
					panic(err)
				}
				err = c.Update(
					bson.M{"_id": domain.Id, "departments._id": id},
					bson.M{"$unset": bson.M{"departments.$.building": ""}})
				if err != nil {
					panic(err)
				}
				return "success"
			}

		})
		r.Post("/:id/category", RequireLogin(), func(domain Domain, u User, req *http.Request, db *mgo.Database, p martini.Params) string {
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

			err := c.Update(bson.M{"_id": domain.Id, "departments._id": id, "departments.category._id": cat}, bson.M{"$pull": bson.M{"departments.$.category": bson.M{"_id": cat}}})
			if err != nil {
				panic(err)
			}
			return "success"
		})
	})
}
