package main

import (
	"github.com/go-martini/martini"

	"encoding/json"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func InitializeDepartmentService(m *martini.ClassicMartini) {
	m.Group("/o/departments", func(r martini.Router) {
		r.Get("/list", RequireLogin(), func(db *mgo.Database) string {
			c := db.C(DepartmentsC)
			var d []Department

			c.Find(bson.M{}).All(&d)
			j, _ := json.Marshal(d)
			return string(j)
		})
		r.Post("/new", RequireLogin(), func(req *http.Request, db *mgo.Database) string {
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

			c := db.C(DepartmentsC)
			err = c.Insert(d)

			if err != nil {
				panic(err)
			}
			return d.Id.Hex()
		})
		r.Post("/:id/new_category", RequireLogin(), func(req *http.Request, db *mgo.Database, p martini.Params) string {
			return ""
		})
	})
}
