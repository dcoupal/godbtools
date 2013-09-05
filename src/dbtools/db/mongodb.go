package db

import (
	//"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

import (
	"dbtools/goext"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDB struct {
	dbtype     string
	Host       string
	Database   string
	Collection string
	mQuery     Doc
	port       int
	query      string
}

func (o *MongoDB) Get(k string) interface{} {
	var v interface{} = nil
	v = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func (o *MongoDB) GetDoc(query string) map[string]interface{} {
	var doc map[string]interface{} = nil
	var mQuery map[string]interface{} = nil
	if query == "" {
		mQuery = o.mQuery
	} else {
		//json.Unmarshal([]byte(query), &mQuery)
		mQuery = Str2mQuery(query)
	}
	session, err := mgo.Dial(o.Host)
	if err != nil {
		fmt.Printf("Can't connect to %s\n", o.Host)
		panic(err)
	}
	defer session.Close()
	coll := session.DB(o.Database).C(o.Collection)
	err = coll.Find(mQuery).One(&doc) //debugger
	if err != nil {
		fmt.Printf("Can't find document with query %v\n", mQuery)
		panic(err)
	}
	return doc
}

// test suite for valjsondb is 7.5 sec before using the iterator, 7.8 with iterator
func (o *MongoDB) GetDocs() <-chan Doc {
	// Connect to MongoDB
	session, err := mgo.Dial(o.Host)
	if err != nil {
		fmt.Printf("Can't connect to %s\n", o.Host)
		panic(err)
	}
	// FIXME - not sure where to put this statement.
	//         here it close too early for the goroutine
	//defer session.Close()

	collCon := session.DB(o.Database).C(o.Collection) //debugger
	iter := collCon.Find(bson.M{}).Iter()

	// TODO - may want to size the channel, so it does not wait to read more from the DB ahead of ti
	ch := make(chan Doc)
	go func() {
		for {
			var aDoc map[string]interface{} //debugger
			if iter.Next(&aDoc) == false {
				break
			}
			ch <- aDoc
		}
		ch <- nil
	}()
	return ch
}

func (o *MongoDB) GetQuery() string {
	return o.query
}

func (o *MongoDB) SetDocProvider(host string, path string) {
	o.Host = host
	if matches, ok := goext.GetParts(path, [][]string{{""}, {"databases", "db"}, {}, {"collections", "c"}, {}}); ok {
		o.Database = matches[2]
		o.Collection = matches[4]
	} else if matches, ok := goext.GetParts(path, [][]string{{""}, {}, {}}); ok {
		o.Database = matches[1]
		o.Collection = matches[2]
	} else { // First item is empty
		panic("Invalid path for a MongoDB provider")
	}
}

func (o *MongoDB) SetQuery(query string) {
	o.query = query
	mQuery := make(map[string]interface{})
	// TODO put the conversion to doc in a library
	q, _ := url.ParseQuery(query)
	for f, _ := range q {
		mQuery[f] = q[f]
	}
	o.mQuery = mQuery
}

func Str2mQuery(query string) map[string]interface{} {
	mQuery := make(map[string]interface{})
	// TODO put the conversion to doc in a library
	q, _ := url.ParseQuery(query)
	for f, _ := range q {
		mQuery[f] = q[f][0]
	}
	return mQuery
}
