package main

// Using this site to generate schemas from examples: http://www.jsonschema.net/#

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
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
	flagset.IntVar(&flags.j, "j", 0, "Parallel factor to validate the documents")
	flagset.StringVar(&flags.query, "query", "{}", "Restrict the validation to documents matching this query")
	flagset.BoolVar(&flags.verbose, "verbose", false, "Show more info")
	flagset.BoolVar(&flags.version, "version", false, "Show the version number")
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func validate(flags *Flags) (int, int) {
	nbDoc := 0
	nbInvalid := 0
	var query map[string]interface{} = nil

	nbWorkers := flags.j
	if nbWorkers == 0 {
		nbWorkers = MaxParallelism()
	}
	runtime.GOMAXPROCS(nbWorkers)

	// queue of documents for the workers
	queueDoc := make(chan map[string]interface{}, 100)
	queueRes := make(chan int, nbWorkers)
	// spawn workers
	for i := 0; i < nbWorkers; i++ {
		go worker(i, queueDoc, queueRes, flags)
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
		queueDoc <- doc
		nbDoc += 1
	}
	// Put a number of stopper in the queue to notify the workers that
	// there is no more documents
	for n := 0; n < nbWorkers; n++ {
		queueDoc <- nil
	}
	// read the results
	var invalid int
	for n := 0; n < nbWorkers; n++ {
		invalid = <-queueRes
		nbInvalid += invalid
	}
	return nbDoc, nbInvalid
}

func validateOneDoc(flags *Flags, schema *gojsonschema.JsonSchemaDocument, doc map[string]interface{}) bool {

	validationResult := schema.Validate(doc)
	fmt.Printf("  item %v, isvalid %v\n", doc["_id"], validationResult.IsValid())
	if validationResult.IsValid() == false {
		fmt.Printf("  %v\n", validationResult.GetErrorMessages())
	}
	return (validationResult.IsValid())
}

func worker(id int, queueDoc chan map[string]interface{}, queueRes chan int, flags *Flags) {

	schema, err := gojsonschema.NewJsonSchemaDocument("file://" + flags.checks)
	if err != nil {
		panic(err.Error())
	}

	nbInvalid := 0
	var doc map[string]interface{}
	for {
		// get work item (pointer) from the queue
		doc = <-queueDoc
		if doc == nil {
			break
		}
		if flags.verbose == true {
			fmt.Printf("worker #%d: item %v\n", id, doc["_id"])
		}
		valid := validateOneDoc(flags, schema, doc)
		if valid == false {
			nbInvalid += 1
		}
	}
	queueRes <- nbInvalid
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
		nbDoc, nbInvalid := validate(flags)
		if flags.verbose == true {
			fmt.Printf("\nValidated %d documents, %d have invalid schemas\n", nbDoc, nbInvalid)
		}
	}
	os.Exit(rc)
}
