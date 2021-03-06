package db

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
)

import (
	"dbtools/goext"
)

type TextFile struct {
	dbtype string
	Host   string
	Path   string
	query  string
}

func (o *TextFile) Get(k string) interface{} {
	var v interface{} = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func (o *TextFile) GetDoc(query string) map[string]interface{} {
	var doc map[string]interface{} = nil
	if query == "" {
		// No query, we expect the file to be the whole document
		doc = goext.ReadJsonFile(o.Path)
	} else {
		// TODO support when many JSON documents can exist in the file
	}
	return doc
}

func (o *TextFile) GetDocs() <-chan Doc {
	ch := make(chan Doc)
	go func() {
		file, _ := os.Open(o.Path)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var aDoc map[string]interface{} //debugger
			jsonLine := scanner.Text()
			json.Unmarshal([]byte(jsonLine), &aDoc)
			ch <- aDoc
		}
		ch <- nil
	}()
	return ch
}

func (o *TextFile) GetQuery() string {
	return o.query
}

func (o *TextFile) SetDocProvider(host string, path string) {
	o.Host = host
	o.Path = path
}

func (o *TextFile) SetQuery(query string) {
	o.query = query
}
