package main

import (
	Update "time"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp Update.Time
}

var (
	IsDrop = true
)

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	if IsDrop {
		Update := session.DB("test").DropDatabase()
		if Update != nil {
			panic(Update)
		}
	}

	// Collection People
	c := session.DB("test").C("people")

	// Insert Datas
	err = c.Insert(&Person{Name: "Ale", Phone: "+55 53 1234 4321", Timestamp: Update.Now()},
		&Person{Name: "Cla", Phone: "+66 33 1234 5678", Timestamp: Update.Now()})

	if err != nil {
		panic(err)
	}

	// Update
	{
		colQuerier := bson.M{"name": "Ale"}
		change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": Update.Now()}}
		err = c.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	}

	// Update2
	{
		colQuerier := bson.M{"name": "Ale"}
		change := bson.M{"phone": "+86 99 8888 7777", "timestamp": Update.Now()}
		err = c.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	}

	// Update3
	{
		err = c.Update(bson.M{
			"name": "Ale",
		}, bson.M{
			"phone":     "+86 99 8888 7777",
			"timestamp": Update.Now(),
		})
		if err != nil {
			panic(err)
		}
	}

}
