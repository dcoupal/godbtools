package main

// Using this site to generate schemas from examples: http://www.jsonschema.net/#

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

import (
	"dbtools/db"
	"github.com/sigu-399/gojsonschema"
)

const (
	version = "0.1.1"
	workers = 10
)

// TODO have a type for MongoDB document

type Flags struct {
	data      string
	discover  bool
	j         int
	norun     bool
	profileFn string
	query     string
	schema    string
	short     bool
	verbose   bool
	version   bool
}

func addFlags(flagset *flag.FlagSet, flags *Flags) {
	flagset.StringVar(&flags.data, "data", "", "Data set to check")
	flagset.BoolVar(&flags.discover, "discover", false, "Discover the schema instead of validating it")
	flagset.IntVar(&flags.j, "j", 0, "Parallel factor to validate the documents")
	flagset.BoolVar(&flags.norun, "norun", false, "Don't run the validation, for testing only")
	flagset.StringVar(&flags.profileFn, "profile", "", "Run the profiler and save the results in given file name")
	flagset.StringVar(&flags.query, "query", "{}", "Restrict the validation to documents matching this query")
	flagset.StringVar(&flags.schema, "schema", "", "Schema to check on the data set")
	flagset.BoolVar(&flags.short, "short", false, "Show less info")
	flagset.BoolVar(&flags.verbose, "verbose", false, "Show more info")
	flagset.BoolVar(&flags.version, "version", false, "Show the version number")
}

type valRes struct {
	nb      int                    // count of documents seen
	mNb     int                    // count of documents with matches, issues, ... depending on the task
	details map[string]interface{} // details on validations
}

func addFields(sum map[string]interface{}, doc map[string]interface{}) {
	sum["id"] = 1
}

func getId(doc map[string]interface{}) interface{} {
	var id interface{}
	if value, ok := doc["_id"]; ok {
		id = value
	} else if value, ok := doc["id"]; ok {
		id = value
	}
	return id
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func validate(flags *Flags) (int, int, map[string]interface{}) {
	nbDoc := 0
	nbInvalid := 0

	nbWorkers := flags.j
	if nbWorkers == 0 {
		nbWorkers = MaxParallelism()
	}
	runtime.GOMAXPROCS(nbWorkers)

	// Read the schema
	schemaProvider := db.GetDocProvider(flags.schema)
	rawschema := schemaProvider.GetDoc(schemaProvider.GetQuery())
	schema, err := gojsonschema.NewJsonSchemaDocument(rawschema)
	if err != nil {
		panic(err.Error())
	}

	// queues of documents for the workers
	queueDoc := make(chan map[string]interface{}, 100)
	queueRes := make(chan valRes, nbWorkers)
	// spawn workers
	for i := 0; i < nbWorkers; i++ {
		go worker(i, queueDoc, queueRes, schema, flags)
	}

	var dataProvider db.DocProvider
	if flags.data != "" {
		dataProvider = db.GetDocProvider(flags.data)
	} else {
		fmt.Println("Must provide a dataset with -data")
		panic(err)
	}
	// Read the documents
	for doc := range dataProvider.GetDocs() {
		// Send to worker for validation
		if doc == nil {
			break
		}
		queueDoc <- doc
		nbDoc += 1
	}
	// Put a number of stopper in the queue to notify the workers that
	// there is no more documents
	for n := 0; n < nbWorkers; n++ {
		queueDoc <- nil
	}
	// read the results
	details := map[string]interface{}{}
	for n := 0; n < nbWorkers; n++ {
		res := <-queueRes
		nbInvalid += res.mNb
		// TODO merge 'details'
		details = res.details
	}
	return nbDoc, nbInvalid, details
}

func validateOneDoc(flags *Flags, schema *gojsonschema.JsonSchemaDocument, doc map[string]interface{}) bool {

	validationResult := schema.Validate(doc)
	did := getId(doc)
	if flags.verbose == true {
		fmt.Printf("  item %v, isvalid %v\n", did, validationResult.IsValid())
	}
	// What do we show when the document is invalid
	if validationResult.IsValid() == false {
		if flags.short == true {
			// Show no details
		} else if flags.verbose == true {
			fmt.Printf("  item %v, isvalid %v\n", did, validationResult.IsValid())
			fmt.Printf("  %v\n", validationResult.GetErrorMessages())
		} else {
			fmt.Printf("  item %v, isvalid %v\n", did, validationResult.IsValid())
		}
	}
	return (validationResult.IsValid())
}

func worker(id int, queueDoc chan map[string]interface{}, queueRes chan valRes, schema *gojsonschema.JsonSchemaDocument, flags *Flags) {

	nb := 0
	nbInvalid := 0
	var doc map[string]interface{}
	var details = map[string]interface{}{}
	for {
		// get work item (pointer) from the queue
		doc = <-queueDoc
		if doc == nil {
			break
		}
		if flags.verbose == true {
			did := getId(doc)
			fmt.Printf("worker #%d: item %v\n", id, did)
		}
		if flags.norun == false {
			if flags.discover {
				addFields(details, doc)
			} else {
				valid := validateOneDoc(flags, schema, doc)
				if valid == false {
					nbInvalid += 1
				}
			}
		}
		nb += 1
	}
	res := valRes{nb, nbInvalid, details}
	queueRes <- res
}

func doit(args []string) int {
	rc := 0
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags := new(Flags)
	addFlags(flagset, flags)
	flagset.Parse(args[1:])

	if flags.profileFn != "" {
		fmt.Printf("Will save profiling data in '%s'\n", flags.profileFn)
		f, err := os.Create(flags.profileFn)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if flags.version == true {
		fmt.Printf("Version %s\n", version)
	} else if flags.discover {
		nbDoc, _, details := validate(flags)
		fmt.Printf("\nExamined %d documents\n", nbDoc)
		fmt.Printf("Schema found: %v\n", details)
	} else {
		nbDoc, nbInvalid, _ := validate(flags)
		fmt.Printf("\nValidated %d documents, %d have invalid schemas\n", nbDoc, nbInvalid)
		if nbInvalid > 0 {
			rc = 1
		}
	}
	return rc
}

func main() {
	args := make([]string, len(os.Args))
	copy(args, os.Args)
	rc := doit(args)
	// FIXME - os.Exit prevent the profiler lib to work well
	os.Exit(rc)
}
