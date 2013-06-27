package main

import (
    //"encoding/json"
   "fmt"
	"os"
)

import (
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

func main() {
    var doc map[string]interface{}
    session, err := mgo.Dial("apollo13.local:27017")  //debugger
    if err != nil {
        fmt.Println("Can't connect to host")
        panic(err)
    }
    defer session.Close()
    collCon := session.DB("log").C("acme_usage")
    iter := collCon.Find(bson.M{}).Iter()
    for iter.Next(&doc) {
        // Send to worker for validation
        fmt.Printf("_id: %s\n", doc["_id"])	//debugger
        fmt.Printf("user: %s\n", doc["user"])	//debugger
    }

	os.Exit(0)
}