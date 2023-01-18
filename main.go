package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ProgName         = "lockbench"
	KeyFileSizeLimit = 8096
)

var (
	version = "undefined"

	showVersion = flag.Bool("version", false, "show program version and exit")
	threads     = flag.Uint("threads", 2, "number of threads to run test")
	iterations  = flag.Uint("iterations", 100000000, "number of iterations in each thread")
)

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "%s\n", ProgName)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	flag.PrintDefaults()
}

func run() int {
	flag.CommandLine.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return 0
	}

	log.Printf("test parameters: threads = %d, iterations = %d", *threads, *iterations)
	startEvent := make(chan struct{})
	var (
		wg  sync.WaitGroup
		mux sync.Mutex
		ctr int
	)
	wg.Add(int(*threads))

	t1 := time.Now()
	for i := uint(0); i < *threads; i++ {
		go func() {
			defer wg.Done()
			<-startEvent
			for i := uint(0); i < *iterations; i++ {
				mux.Lock()
				ctr++
				mux.Unlock()
			}
		}()
	}

	close(startEvent)
	wg.Wait()
	t2 := time.Now()
	deltaT := t2.Sub(t1)
	log.Printf("test duration: %s, ctr value = %d", deltaT, ctr)

	return 0
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Default().SetPrefix(strings.ToUpper(ProgName) + ": ")
	os.Exit(run())
}
