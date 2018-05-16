package usesafemongo

import (
	"github.com/podhmo/safemongo/experiment/safemongo"
	bson "gopkg.in/mgo.v2/bson"
)

func f1(c *safemongo.Collection) error {
	return c.Update(bson.M{"_id": bson.NewObjectId()}, bson.M{"name": "foo"})
}

func f2(c *safemongo.Collection) error {
	return c.Update(bson.M{"_id": bson.NewObjectId()}, bson.M{"$set": bson.M{"name": "foo"}})
}

func f3(c *safemongo.Collection) error {
	sv := bson.M{"name": "foo"}
	return c.Update(bson.M{"_id": bson.NewObjectId()}, bson.M{"$set": sv})
}

func f4(c *safemongo.Collection) error {
	sv := bson.M{"name": "foo"}
	sv = bson.M{"$set": sv}
	return c.Update(bson.M{"_id": bson.NewObjectId()}, sv)
}

func f5(c *safemongo.Collection) error {
	sv := func(sv bson.M) bson.M { return bson.M{"$set": sv} }
	return c.Update(bson.M{"_id": bson.NewObjectId()}, sv(bson.M{"name": "foo"}))
}
