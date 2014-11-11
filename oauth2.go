package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/oauth2"
	"github.com/golang/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
)

func ConfigureOauth2() {

	cfg, err := google.NewConfig(&oauth2.Options{
		ClientID:     "812975936151-i8po4eflb6fggohokgl98d5998uh4t6k.apps.googleusercontent.com",
		ClientSecret: "5oqoK2q-_lnHO5kCdB8DjSyh",
		RedirectURL:  "http://localhost/o/token",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	})
	if err != nil {
		print("Uh Oh!")
	}
	Cfg = cfg
}

func Oauth2Token(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	code := u.Query().Get("code")
	if code != "" {
		transport, err := Cfg.NewTransportWithCode(code)
		if err != nil {
			log.Print(err)
			return
		}
		client := http.Client{Transport: transport}
		resp, err := client.Get("https://www.googleapis.com/userinfo/v2/me")

		defer resp.Body.Close()

		if err != nil {
			//log.Fatal(err)
			return
		}
		gData := new(GoogleUserV2)

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return
		}

		err = json.Unmarshal([]byte(contents), &gData)

		if err != nil {
			print(err.Error())
			return
		}

		print("User Exists:")
		exists, user, err := UserExists(MongoSession, gData.Id)
		if err != nil {
			println("ERROR - " + err.Error())
			return
		}
		println(exists)

		if !exists {

			print("UserCreate:")
			newUser, err := UserCreate(MongoSession, *gData)

			if err != nil {
				println("ERROR - " + err.Error())
				return
			}
			println(newUser.Id)
			out, err := newUser.Marshall()
			fmt.Fprint(w, string(out))
		} else {
			out, err := user.Marshall()
			if err == nil {
				fmt.Fprint(w, string(out))
			}
		}

	} else {
		print("Uh oh!")
	}
}

func Oauth2Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, Cfg.AuthCodeURL("token", "offline", "force"), http.StatusFound)
}
