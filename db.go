package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var Mgo *mgo.Session

const (
	COLLECTION_NAME_MESSAGES = "messages"
)

func init() {
	session, err := mgo.Dial("mongodb://localhost/chat")
	if err != nil {
		panic(err)
	}

	Mgo = session
}

func save(msg *Message) {
	msg.Timestamp = time.Now()
	if err := Mgo.DB("").C(COLLECTION_NAME_MESSAGES).Insert(msg); err != nil {
		log.Print(err)
	}
}
