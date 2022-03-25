package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Mgo *mgo.Session

const (
	DEFAULT_LIMIT            = 100
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

func find(from string, to string, limit int64) ([]Message, error) {
	var findQuery []bson.M
	var query []bson.M
	coll := Mgo.DB("").C(COLLECTION_NAME_MESSAGES)

	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	if to != "" {
		query = append(query, bson.M{"to": to})
	}

	if from != "" {
		query = append(query, bson.M{"from": from})
	}

	// By default, only messages from the last 30 days should be returned.
	query = append(query, bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -30)}})

	findQuery = append(findQuery, bson.M{"$match": func() bson.M {
		if query != nil && len(query) > 0 {
			return bson.M{"$and": query}
		}
		return bson.M{}
	}()})

	findQuery = append(findQuery, bson.M{
		"$sort": bson.M{"timestamp": -1},
	}, bson.M{
		"$limit": limit,
	})

	messages := make([]Message, limit)
	log.Print("query: ", findQuery)
	err := coll.Pipe(findQuery).All(&messages)

	if err != nil {
		return nil, err
	}

	return messages, nil
}
