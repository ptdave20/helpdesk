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
			//err = c.UpdateId(domain.Id, bson.M{"departments": bson.M{"$addToSet": d}})
			//err = c.Update(bson.M{"_id": domain.Id}, bson.M{"departments": bson.M{"$push": d}})
			if err != nil {
				panic(err)
			}
			return d.Id.Hex()
		})
		r.Post("/:id", RequireLogin(), func(req *http.Request, db *mgo.Database, p martini.Params) string {
			return ""
		})
	})
}
