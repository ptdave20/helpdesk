package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var CFG Config

func (cfg Config) Load() bool {
	cfgBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return false
	}
	json.Unmarshal(cfgBytes, &cfg)

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

	CFG.ClientID = clientid
	CFG.ClientSecret = secret
	CFG.RedirectURI = host + "/o/token"
	CFG.Hostname = host
	CFG.MongoAddress = mongo
	CFG.MongoDatabase = db

	return true
}

func (cfg Config) Save() (bool, error) {
	cfgBytes, err := json.Marshal(&cfg)
	err = ioutil.WriteFile("config.json", cfgBytes, 0755)
	if err != nil {
		return false, err
	}
	return true, nil
}
