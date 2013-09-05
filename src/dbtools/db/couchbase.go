package db

import (
	"fmt"
	"reflect"
)

import (
	"dbtools/goext"
	"github.com/couchbaselabs/go-couchbase"
)

type CouchBase struct {
	dbtype string
	Design string
	Host   string
	Bucket string
	Pool   string
	port   int
	Query  string
	View   string
}

func (o *CouchBase) Get(k string) interface{} {
	var v interface{} = reflect.ValueOf(o).Elem().FieldByName(k).Interface()
	return v
}

func getLocalBucket(bucketN string) *couchbase.Bucket {
	c, err := couchbase.Connect("http://" + "localhost:8091")
	mf(err, "connect - "+"localhost:8091")
	p, err := c.GetPool("default")
	mf(err, "pool")
	b, err := p.GetBucket(bucketN)
	mf(err, "bucket")
	defer b.Close()
	mf(err, "bucket.get")
	return b
}

func (o *CouchBase) GetDoc(query string) map[string]interface{} {
	var doc map[string]interface{} = nil

	// TODO - Code is from Couchbase example, clean it to make it more uniform to current code base
	c, err := couchbase.Connect("http://" + o.Host)
	mf(err, "connect - "+o.Host)
	p, err := c.GetPool(o.Pool)
	mf(err, "pool")
	b, err := p.GetBucket(o.Bucket)
	mf(err, "bucket")
	defer b.Close()
	err = b.Get(o.Query, &doc)
	mf(err, "bucket.get")

	return doc
}

func (o *CouchBase) GetDocs() <-chan Doc {
	c, err := couchbase.Connect("http://" + o.Host)
	mf(err, "connect - "+o.Host)
	p, err := c.GetPool(o.Pool)
	mf(err, "pool")
	b, err := p.GetBucket(o.Bucket)
	mf(err, "bucket")
	// FIXME - not sure where to put this statement.
	//         here it close too early for the goroutine
	//defer b.Close()

	// TODO - may want to size the channel, so it does not wait to read more from the DB ahead of time
	ch := make(chan Doc)
	go func() {
		vres, _ := b.View(o.Design, o.View, map[string]interface{}{})
		var aDoc map[string]interface{}
		for i := 0; i < vres.TotalRows; i++ {
			cbRow := vres.Rows[i]
			cbValue := cbRow.Value
			aDoc = getMap(cbValue) //debugger
			ch <- aDoc
		}
		ch <- nil
		b.Close()
	}()
	return ch
}

func (o *CouchBase) GetQuery() string {
	return o.Query
}

func getMap(value interface{}) map[string]interface{} {
	res := map[string]interface{}{} //debugger
	for k, v := range value.(map[string]interface{}) {
		res[k] = v
	}
	return res
}

func getMap2(value interface{}) map[string]interface{} {
	var res map[string]interface{} //debugger
	switch value.(type) {
	case map[interface{}]interface{}:
		for k, v := range value.(map[interface{}]interface{}) {
			if str, ok := k.(string); ok {
				res[str] = v
			}
		}
	default:
		fmt.Printf("ERROR - can't convert to map[string]interface{}\n%v\n\n", value)
	}
	return res
}

func (o *CouchBase) SetDocProvider(host string, path string) {
	if o.Host != "" {
		o.Host = host
	} else {
		o.Host = "localhost:8091"
	}
	o.Pool = "default" //debugger
	if matches, ok := goext.GetParts(path, [][]string{{""}, {"buckets"}, {}, {"_design", "designs"}, {}, {"_view", "views"}, {}}); ok {
		o.Bucket = matches[2]
		o.Design = matches[4]
		o.View = matches[6]
	} else if matches, ok := goext.GetParts(path, [][]string{{""}, {}, {"_design", "designs"}, {}, {"_view", "views"}, {}}); ok {
		o.Bucket = matches[1]
		o.Design = matches[3]
		o.View = matches[5]
	} else if matches, ok := goext.GetParts(path, [][]string{{""}, {"buckets"}, {}}); ok {
		o.Bucket = matches[2]
	} else if matches, ok := goext.GetParts(path, [][]string{{""}, {}}); ok {
		o.Bucket = matches[1]
	} else if matches, ok := goext.GetParts(path, [][]string{{""}, {}, {}}); ok {
		o.Pool = matches[1]
		o.Bucket = matches[2]
	} else {
		panic("Invalid path for a CouchBase provider: " + path)
	}
}

func (o *CouchBase) SetQuery(query string) {
	o.Query = query
}

func mf(err error, msg string) {
	if err != nil {
		//log.Fatalf("%v: %v", msg, err)
		fmt.Printf("%v", msg)
		panic(err)
	}
}
