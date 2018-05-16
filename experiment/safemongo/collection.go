package safemongo

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

// Collection :
type Collection struct {
	*mgo.Collection
}

// Insert is mgo.Collection.Insert wrapper function.
func (c Collection) Insert(docs ...interface{}) error {
	now := timeNow()
	for i, d := range docs {
		docs[i] = withUpdatedAt(d, now)
	}
	return c.Collection.Insert(docs...)
}

// Update is mgo.Collection.Update wrapper function.
func (c Collection) Update(selector interface{}, update interface{}) error {
	return c.Collection.Update(selector, withUpdatedAtForUpdate(update))
}

// UpdateId is mgo.Collection.UpdateId wrapper function.
func (c Collection) UpdateId(id interface{}, update interface{}) error {
	return c.Update(bson.D{{Name: "_id", Value: id}}, update)
}

// UpdateAll is mgo.Collection.UpdateAll wrapper function.
func (c Collection) UpdateAll(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Collection.UpdateAll(selector, withUpdatedAtForUpdate(update))
}

// Upsert is mgo.Collection.Upsert wrapper function.
func (c Collection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Collection.Upsert(selector, withUpdatedAt(update, timeNow()))
}

// UpsertId is mgo.Collection.UpsertId wrapper function.
func (c Collection) UpsertId(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.Upsert(bson.D{{Name: "_id", Value: id}}, update)
}

var updatedAtFieldName = "_updatedAt"
var updatedAtFieldTag = fmt.Sprintf(`bson:"%s"`, updatedAtFieldName)

func withUpdatedAt(d interface{}, updatedAt time.Time) interface{} {
	if d == nil {
		return nil
	}

	v := reflect.ValueOf(d)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if m, ok := v.Interface().(bson.M); ok {
		if isBsonMWithOperator(m) {
			if _, ok := m["$set"]; ok {
				m["$set"] = withUpdatedAt(m["$set"], updatedAt)
			} else {
				m["$set"] = bson.M{updatedAtFieldName: updatedAt}
			}
			return m
		}
	}

	tp := reflect.StructOf([]reflect.StructField{
		{Name: "D", Type: v.Type(), Tag: reflect.StructTag(`bson:",inline"`)},
		{Name: "UpdatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(updatedAtFieldTag)},
	})
	e := reflect.New(tp).Elem()
	e.FieldByName("D").Set(v)
	e.FieldByName("UpdatedAt").Set(reflect.ValueOf(updatedAt))
	return e.Interface()
}

func withUpdatedAtForUpdate(update interface{}) interface{} {
	switch m := update.(type) {
	case bson.M:
		if !isBsonMWithOperator(m) {
			panic(fmt.Sprintf("%s", update))
		}
		update = withUpdatedAt(update, timeNow())
	default:
		panic(fmt.Sprintf("%s", update))
	}
	return update
}

func isBsonMWithOperator(m bson.M) bool {
	for k := range m {
		if strings.HasPrefix(k, "$") {
			return true
		}
	}
	return false
}
