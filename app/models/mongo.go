package models

import (
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
)

var mongo *mgo.Session = nil

func DB() (*mgo.Database, error) {
	if mongo == nil {
		var err error
		mongo, err = mgo.Dial("localhost")
		if err != nil {
			return nil, err
		}
		mongo.SetMode(mgo.Monotonic, true)
	}
	return mongo.DB("byvnotes"), nil
}
