package main

import (
	"github.com/golang/oauth2"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"sync"
)

var Cfg *oauth2.Config
var MongoSession *mgo.Session
var Sessions = sessions.NewCookieStore([]byte("helpdesk"))

func main() {
	ConfigureOauth2()
	session, err := mgo.Dial("localhost")

	if err != nil {
		log.Fatal(err)
	} else {
		MongoSession = session
	}

	http.HandleFunc("/o/login", Oauth2Login)
	http.HandleFunc("/o/token", Oauth2Token)
	http.ListenAndServe(":80", nil)
}
