package db

import (
	"dbtools/goext"
	"reflect"
)

type TextFile struct {
	dbtype string
	host   string
	path   string
	query  string
}

func (o *TextFile) get(k string) interface{} {
	var v interface{} = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func (o *TextFile) GetDoc(query string) map[string]interface{} {
	var doc map[string]interface{} = nil
	if query == "" {
		// No query, we expect the file to be the whole document
		doc = goext.ReadJsonFile(o.path) //debugger
	} else {
		// TODO support when many JSON documents can exist in the file
	}
	return doc
}

func (o *TextFile) GetQuery() string {
	return o.query
}

func (o *TextFile) SetDocProvider(host string, path string) {
	o.host = host
	o.path = path
}

func (o *TextFile) SetQuery(query string) {
	o.query = query
}
