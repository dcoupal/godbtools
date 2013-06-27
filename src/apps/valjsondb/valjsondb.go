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

type Flags struct {
	checks     string
	collection string
	connection string
	database   string
	query      string
	verbose    bool
	version    bool
}

func addFlags(flagset *flag.FlagSet, flags *Flags) {
	flagset.StringVar(&flags.checks, "checks", "", "Checks to run on the collection")
	flagset.StringVar(&flags.collection, "collection", "", "Collection to validate")
	flagset.StringVar(&flags.connection, "connection", "localhost:27017", "Connection to the database, if none try locally")
	flagset.StringVar(&flags.database, "database", "", "Database to check")
	flagset.StringVar(&flags.query, "query", "{}", "Restrict the validation to documents matching this query")
	flagset.BoolVar(&flags.verbose, "verbose", false, "Show more info")
	flagset.BoolVar(&flags.version, "version", false, "Show the version number")
}

func validate(flags *Flags) {
	var doc map[string]interface{}
	var query map[string]interface{} = nil
	// Connect to the DB

	// Read the documents
	session, err := mgo.Dial(flags.connection)
	if err != nil {
		fmt.Printf("Can't connect to %s\n", flags.connection)
		panic(err)
	}
	defer session.Close()
	collCon := session.DB(flags.database).C(flags.collection)
	json.Unmarshal([]byte(flags.query), &query)
	iter := collCon.Find(query).Iter()
	for iter.Next(&doc) {
		// Send to worker for validation
		validateOneDoc(flags, doc)
	}
}

func validateOneDoc(flags *Flags, doc map[string]interface{}) {
	if flags.verbose == true {
		fmt.Printf("_id: %s\n", doc["_id"])
	}
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
	fmt.Printf("IsValid %v\n", validationResult.IsValid())
	fmt.Printf("%v\n", validationResult.GetErrorMessages())
}

func main() {
	rc := 0
	args := make([]string, len(os.Args))
	copy(args, os.Args)
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError) //debugger
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
