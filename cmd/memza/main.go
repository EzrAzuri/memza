package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rcompos/memza/memza"
)

func main() {

	var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB
	var filePut, fileGet, fileOut, memcachedServer string
	var test, debug, force bool
	flag.StringVar(&filePut, "p", "", "File to put")
	flag.StringVar(&fileGet, "g", "", "File to get")
	flag.StringVar(&fileOut, "o", "out.dat", "Output file for retrieval")
	flag.StringVar(&memcachedServer, "s", "localhost:11211", "memcached_server:port")
	flag.BoolVar(&debug, "d", false, "Debug mode")
	flag.BoolVar(&test, "t", false, "Check memcached server")
	flag.BoolVar(&force, "f", false, "Force key overwrite")
	flag.Parse()

	//fmt.Printf("flags> %s, %s, %s\n", filePut, fileGet, fileOut)
	if test == true && (filePut != "" || fileGet != "") {
		memza.HelpMe("Must supply test as single argument (-t).")
	}

	if filePut != "" && fileGet != "" {
		memza.HelpMe("Must supply file as argument (-p or -g).")
	}

	if filePut == "" && fileGet == "" && test == false {
		memza.HelpMe("")
	}

	if filePut != "" {
		if err := memza.StoreFile(filePut, memcachedServer, maxFileSize, debug, force); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if fileGet != "" {
		if err := memza.RetrieveFile(fileGet, memcachedServer, fileOut, debug); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if test == true {
		if err := memza.CheckServer(memcachedServer); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

}
