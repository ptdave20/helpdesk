package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func GetDomain(mongo *mgo.Session, string domain) (*Domain, err) {
	session := mongo.Clone()

	defer session.Close()

	col := session.DB("helpdesk").C("domain")
	var domains []Domain
	err := col.Find(bson.M{"domain": domain}).All(&domains)

	if err != nil {
		return nil, err
	}

	if len(domains) > 0 {
		return &domains[0], nil
	}
}
