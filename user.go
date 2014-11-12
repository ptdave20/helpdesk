package main

import (
	//"encoding/json"
	//"github.com/go-martini/martini"
	//"github.com/golang/oauth2"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

func UserExists(mSession *mgo.Session, Id string) (bool, *User, error) {

	session := mSession.Clone()

	// Close our session
	defer session.Close()

	var users []User

	col := session.DB("helpdesk").C("user")
	err := col.Find(bson.M{"google_id": Id}).All(&users)
	if err != nil {
		return false, new(User), err
	}
	if len(users) == 0 {
		return false, nil, nil
	}

	return true, &users[0], nil
}

func UserCreate(mSession *mgo.Session, user GoogleUserV2) (User, error) {
	session := mSession.Clone()

	// Close our session and free our sync
	defer session.Close()

	var newUser User
	newUser.Id = bson.NewObjectId()
	newUser.Firstname = user.GivenName
	newUser.Lastname = user.FamilyName
	newUser.Email = user.Email
	newUser.GoogleId = user.Id

	newUser.FirstLogin = time.Now()
	newUser.LastLogin = time.Now()
	newUser.Enabled = true

	col := session.DB("helpdesk").C("user")
	err := col.Insert(newUser)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}
