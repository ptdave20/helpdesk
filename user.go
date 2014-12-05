package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/golang/oauth2"
	"github.com/golang/oauth2/google"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
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
	newUser.Picture = gUser.Picture
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

func RequireLogin() martini.Handler {
	return func(db *mgo.Database, s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		users := db.C(UsersC)
		user_session := db.C(SessionsC)
		if s.Get("session") != nil {
			session_id := s.Get("session")
			u := new(User)
			ses := new(Session)
			err := user_session.Find(bson.M{"_id": bson.ObjectIdHex(session_id.(string))}).One(ses)
			if err != nil {
				http.Redirect(w, r, "/o/login?state="+r.URL.Path, 302)
			}
			err = users.FindId(ses.UserId).One(u)
			if err != nil {
				http.Redirect(w, r, "/o/login?state="+r.URL.Path, 302)
			}
			if !ses.Expired() {
				c.Map(*u)
				c.Map(*ses)
				c.Next()
			} else {
				http.Redirect(w, r, "/o/login?state="+r.URL.Path, 302)
			}

		} else {
			http.Redirect(w, r, "/o/login?state="+r.URL.Path, 302)
		}
	}
}

func RequireLoginNoRedirectOutResult() martini.Handler {
	return func(db *mgo.Database, s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		users := db.C(UsersC)
		user_session := db.C(SessionsC)
		var res SimpleResult
		if s.Get("session") != nil {
			session_id := s.Get("session")
			u := new(User)
			ses := new(Session)
			err := user_session.Find(bson.M{"_id": bson.ObjectIdHex(session_id.(string))}).One(ses)
			if err != nil {
				res.Result = false
			} else {
				err = users.FindId(ses.UserId).One(u)
				if err != nil {
					res.Result = false
				} else {
					if ses.Expired() {
						res.Result = false
					} else {
						res.Result = true
					}
				}
			}
		} else {
			res.Result = false
		}
		b, _ := json.Marshal(res)
		fmt.Fprint(w, string(b))
	}
}

func Oauth2Handler() martini.Handler {

	cfg, _ := oauth2.New(
		oauth2.Client(CFG.ClientID, CFG.ClientSecret),
		oauth2.RedirectURL(CFG.RedirectURI),
		oauth2.Scope("https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"),
		google.Endpoint(),
	)

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
				transport, err := cfg.NewTransportFromCode(code)
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
				if next == "" {
					next = "/"
				}
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
	var ns = Session{bson.NewObjectId(), user.Id, *token}
	ns.Id = bson.NewObjectId()

	col.Insert(ns)

	return ns.Id.Hex()
}

func GetUserById(db *mgo.Database, id string) *User {
	u := new(User)
	err := db.C(UsersC).Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(u)
	if err != nil {
		return nil
	}
	return u
}

func InitializeUserService(m *martini.ClassicMartini) {
	m.Group("/o/user", func(r martini.Router) {
		r.Get("/logged_in", RequireLoginNoRedirectOutResult())
		r.Get("/list", RequireLogin(), func(db *mgo.Database) string {
			var users []User
			c := db.C(UsersC)
			c.Find(bson.M{}).All(&users)
			b, _ := json.Marshal(&users)
			return string(b)
		})
		r.Get("/me", RequireLogin(), func(u User) string {
			b, err := u.Marshal()
			if err != nil {
				panic(err)
			}
			return string(b)
		})
		r.Get("/:id", RequireLogin(), func(db *mgo.Database, p martini.Params) string {
			u := GetUserById(db, p["id"])
			b, _ := u.Marshal()
			return string(b)
		})

	})
}
