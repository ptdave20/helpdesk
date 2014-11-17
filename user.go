package main

import (
<<<<<<< HEAD
	//"encoding/json"
	//"github.com/go-martini/martini"
	//"github.com/golang/oauth2"
=======
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/golang/oauth2"
	"github.com/golang/oauth2/google"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
>>>>>>> 52179a7d6a38c1c9f904c475275d43bf740a9667
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

func UserExists(db *mgo.Database, Id string) (*User, error) {
	user := new(User)

	col := db.C(UsersC)
	err := col.Find(bson.M{"google_id": Id}).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UserCreate(db *mgo.Database, gUser GoogleUserV2) (*User, error) {
	newUser := new(User)
	newUser.Id = bson.NewObjectId()
	newUser.Firstname = gUser.GivenName
	newUser.Lastname = gUser.FamilyName
	newUser.Email = gUser.Email
	newUser.GoogleId = gUser.Id

	newUser.FirstLogin = time.Now()
	newUser.LastLogin = time.Now()
	newUser.Enabled = true

	col := db.C(UsersC)
	err := col.Insert(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
<<<<<<< HEAD
=======

func RequireLogin() martini.Handler {
	return func(db *mgo.Database, s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		users := db.C(UsersC)
		user_session := db.C(SessionsC)
		if s.Get("session") != nil {
			session_id := s.Get("session")
			print(session_id.(string))
			u := new(User)
			ses := new(Session)
			err := user_session.Find(bson.M{"_id": bson.ObjectIdHex(session_id.(string))}).One(ses)
			if err != nil {
				panic(err)
				//c.Next()
			}
			err = users.FindId(ses.UserId).One(u)
			if err != nil {
				panic(err)
				//c.Next()
			}
			c.Map(*u)
			c.Map(*ses)
			c.Next()
		} else {
			http.Redirect(w, r, "/o/login?state="+r.URL.Path, 302)
		}
	}
}

func Oauth2Handler() martini.Handler {
	cfg, _ := google.NewConfig(&oauth2.Options{
		ClientID:     "812975936151-i8po4eflb6fggohokgl98d5998uh4t6k.apps.googleusercontent.com",
		ClientSecret: "5oqoK2q-_lnHO5kCdB8DjSyh",
		RedirectURL:  "http://localhost/o/token",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	})

	return func(db *mgo.Database, s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/o/login":
				state := r.URL.Query().Get("state")
				if state == "" {
					state = "/"
				}
				http.Redirect(w, r, cfg.AuthCodeURL(state, "offline", "force"), 302)
				break
			case "/o/token":
				u := r.URL
				next := u.Query().Get("state")
				code := u.Query().Get("code")
				transport, err := cfg.NewTransportWithCode(code)
				if err != nil {
					print(err.Error())
					//http.Redirect(w, r, "/o/error", 302)
					return
				}
				gUser := GetGoogleUser(transport)
				user := FindOrCreateUser(db, gUser)
				ses := CreateUserSession(db, *user, transport.Token())
				s.Set("session", ses)
				//val, _ := json.Marshal(transport.Token())
				//s.Set("helpdesk", val)
				http.Redirect(w, r, next, 302)
				break
			case "/o/logout":
				u := r.URL
				next := u.Query().Get("next")
				s.Clear()
				http.Redirect(w, r, next, 302)
				break
			}
		}
		c.Map(cfg)
	}
}

func GetGoogleUser(transport *oauth2.Transport) GoogleUserV2 {
	client := http.Client{Transport: transport}
	resp, _ := client.Get("https://www.googleapis.com/userinfo/v2/me")
	defer resp.Body.Close()
	g := new(GoogleUserV2)
	contents, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(contents, g)
	if err != nil {
		panic(err)
	}
	return *g
}

func FindOrCreateUser(db *mgo.Database, google_user GoogleUserV2) *User {
	user, err := UserExists(db, google_user.Id)

	if user != nil {
		return user
	}

	user, err = UserCreate(db, google_user)
	if err != nil {
		panic(err)
	}
	return user
}

func CreateUserSession(db *mgo.Database, user User, token *oauth2.Token) string {
	col := db.C(SessionsC)
	ns := new(Session)
	ns.Active = true
	ns.Expires = token.Expiry
	ns.Refresh = token.RefreshToken
	ns.Token = token.AccessToken
	ns.UserId = user.Id
	ns.Id = bson.NewObjectId()

	col.Insert(ns)

	return ns.Id.Hex()
}
>>>>>>> 52179a7d6a38c1c9f904c475275d43bf740a9667
