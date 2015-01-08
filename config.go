package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func (cfg Config) Load() bool {
	cfgBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return false
	}
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return false
	}

	return true
}

func (cfg Config) FirstTimeSetup() bool {
	print("Enter your ClientID: ")
	var clientid, secret, host, mongo, db string
	_, err := fmt.Scanf("%s\n", &clientid)
	if err != nil {
		return false
	}

	print("Enter your Client Secret: ")
	_, err = fmt.Scanf("%s\n", &secret)
	if err != nil {
		return false
	}

	print("Enter your Host Name: http://")
	_, err = fmt.Scanf("%s\n", &host)
	if err != nil {
		return false
	}

	print("Enter your MongoDB Address: mongodb://")
	_, err = fmt.Scanf("%s\n", &mongo)
	if err != nil {
		return false
	}

	print("Enter your MongoDB Database: ")
	_, err = fmt.Scanf("%s\n", &db)
	if err != nil {
		return false
	}

	cfg.ClientID = clientid
	cfg.ClientSecret = secret
	cfg.RedirectURI = host + "/o/token"
	cfg.Hostname = host
	cfg.MongoAddress = mongo
	cfg.MongoDatabase = db

	cfg.SetupDB()

	return true
}

func (cfg Config) SetupDB() (bool, error) {
	session, err := mgo.Dial(cfg.MongoAddress)
	if err != nil {
		return false, err
	}
	defer session.Close()

	db := session.DB(cfg.MongoDatabase)

	// Check for domains
	dom := db.C(DomainsC)
	count, err := dom.Find(bson.M{}).Count()
	if err != nil {
		return false, nil
	}

	// We need to create our initial domain
	if count == 0 {
		var domain Domain

		print("Enter the name of your domain (Your company): ")
		_, err = fmt.Scanf("%s\n", &domain.Name)
		if err != nil {
			return false, err
		}

		print("Enter the domain name (yourcompany.org): ")
		var tmp string
		_, err = fmt.Scanf("%s\n", &tmp)
		if err != nil {
			return false, err
		}
		domain.AcceptedDomains = append(domain.AcceptedDomains, tmp)

		err = dom.Insert(domain)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return true, nil
}

func (cfg Config) Save() (bool, error) {
	cfgBytes, err := json.Marshal(&cfg)
	err = ioutil.WriteFile("config.json", cfgBytes, 0755)
	if err != nil {
		return false, err
	}
	return true, nil
}
