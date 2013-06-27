package main

// Using this site to generate schemas from examples: http://www.jsonschema.net/#

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

import (
	"github.com/sigu-399/gojsonschema"
	"labix.org/v2/mgo"
)

const (
	version = "0.1.0"
	workers = 10
)

// TODO have a type for MongoDB document

type Flags struct {
	checks     string
	collection string
	connection string
	database   string
	j          int
	query      string
	verbose    bool
	version    bool
}

func addFlags(flagset *flag.FlagSet, flags *Flags) {
	flagset.StringVar(&flags.checks, "checks", "", "Checks to run on the collection")
	flagset.StringVar(&flags.collection, "collection", "", "Collection to validate")
	flagset.StringVar(&flags.connection, "connection", "localhost:27017", "Connection to the database, if none try locally")
	flagset.StringVar(&flags.database, "database", "", "Database to check")
	flagset.IntVar(&flags.j, "j", 1, "Parallel factor to validate the documents")
	flagset.StringVar(&flags.query, "query", "{}", "Restrict the validation to documents matching this query")
	flagset.BoolVar(&flags.verbose, "verbose", false, "Show more info")
	flagset.BoolVar(&flags.version, "version", false, "Show the version number")
}

func validate(flags *Flags) {
	var query map[string]interface{} = nil

	// queue of documents for the workers
	queue := make(chan map[string]interface{})
	// spawn workers
	for i := 0; i < flags.j; i++ {
		go worker(i, queue, flags)
	}

	// Connect to the DB
	session, err := mgo.Dial(flags.connection)
	if err != nil {
		fmt.Printf("Can't connect to %s\n", flags.connection)
		panic(err)
	}
	defer session.Close()
	collCon := session.DB(flags.database).C(flags.collection)

	// Read the documents
	json.Unmarshal([]byte(flags.query), &query)
	iter := collCon.Find(query).Iter()
	for {
		var doc map[string]interface{}
		if iter.Next(&doc) == false {
			break
		}
		// Send to worker for validation
		queue <- doc
	}
	// Put a number of stopper in the queue to notify the workers that
	// there is no more documents
	for n := 0; n < flags.j; n++ {
		queue <- nil
	}
}

func validateOneDoc(flags *Flags, doc map[string]interface{}) {
	//schema, err := gojsonschema.NewJsonSchemaDocument("http://myhost/bla/schema1.json")
	// OR
	schema, err := gojsonschema.NewJsonSchemaDocument("file://" + flags.checks)

	if err != nil {
		panic(err.Error())
	}

	//jsonToValidate, err := json.Marshal(doc)
	//if err != nil {
	//	panic(err.Error())
	//}
	validationResult := schema.Validate(doc)
	fmt.Printf("item %v\n", doc["_id"])
	fmt.Printf("IsValid %v\n", validationResult.IsValid())
	fmt.Printf("%v\n", validationResult.GetErrorMessages())
}

func worker(id int, queue chan map[string]interface{}, flags *Flags) {
	var doc map[string]interface{}
	for {
		// get work item (pointer) from the queue
		doc = <-queue
		if doc == nil {
			break
		}
		if flags.verbose == true {
			fmt.Printf("worker #%d: item %v\n", id, doc["_id"])
		}
		validateOneDoc(flags, doc)
	}
}

func main() {
	rc := 0
	args := make([]string, len(os.Args))
	copy(args, os.Args)
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags := new(Flags)
	addFlags(flagset, flags)
	flagset.Parse(args[1:])
	if flags.version == true {
		fmt.Printf("Version %s\n", version)
	} else {
		validate(flags)
	}
	os.Exit(rc)
}
