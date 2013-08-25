package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		prof, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err.Error())
		}
		defer prof.Close()

		pprof.StartCPUProfile(prof)
		defer pprof.StopCPUProfile()
	}

	go func() {
		x := make(chan bool, 1)
		for {
			x <- true
			<-x
			log.Printf("boo")
			runtime.Gosched()
		}
	}()

	// Waiting for terminating (i use a sighandler like in vitess)
	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)
	<-terminate

	log.Printf("Server stopped")
}
