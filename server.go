package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	//"github.com/martini-contrib/gzip"
	"fmt"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
	"log"
)

var CFG Config

func main() {
	println("Looking for config.json")
	cfgBytes, err := ioutil.ReadFile("config.json")

	// There was an error reading the config file, may not exist
	if err != nil {
		// Attempt to read the data in from the console

		print("Enter your ClientID     : ")
		var clientid, secret, host string
		_, err := fmt.Scanf("%s\n", &clientid)
		if err != nil {
			log.Fatal(err)
		}

		print("Enter your Client Secret: ")
		_, err = fmt.Scanf("%s\n", &secret)
		if err != nil {
			log.Fatal(err)
		}

		print("Enter your Host Name    : ")
		_, err = fmt.Scanf("%s\n", &host)
		if err != nil {
			log.Fatal(err)
		}

		CFG.ClientID = clientid
		CFG.ClientSecret = secret
		CFG.RedirectURI = host + "/o/token"

		cfgBytes, err = json.Marshal(&CFG)

		err = ioutil.WriteFile("config.json", cfgBytes, 0755)

		// We were unable to save, continue on anyhow
		if err != nil {
			println("Unable to save config")
		}
	}
	err = json.Unmarshal(cfgBytes, &CFG)

	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()

	m.Use(martini.Static("assets", martini.StaticOptions{Fallback: "/index.html", Exclude: "/o"}))
	m.Use(MongoDB())

	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("session"))))
	m.Use(Oauth2Handler())

	InitTicketService(m)
	InitializeUserService(m)
	InitializeDepartmentService(m)
	InitTicketService(m)

	//m.Use(gzip.All())
	m.RunOnAddr(":80")

}
