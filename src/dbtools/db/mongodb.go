package db

import (
	//"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

import (
	"labix.org/v2/mgo"
)

type MongoDB struct {
	dbtype     string
	Host       string
	port       int
	Database   string
	Collection string
	query      string
	mQuery     Doc
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

func (o *MongoDB) GetQuery() string {
	return o.query
}

func (o *MongoDB) SetDocProvider(host string, path string) {
	o.Host = host
	arr := strings.Split(path, "/")
	// TODO - test that we have enough fields
	o.Database = arr[1]
	o.Collection = arr[2]
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
