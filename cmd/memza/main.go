package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rcompos/memza/memza"
)

func main() {

	var memcachedServer string = "localhost:11211"
	var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB
	var filePut, fileGet, fileOut string
	var check, debug bool
	flag.StringVar(&filePut, "p", "", "File to put")
	flag.StringVar(&fileGet, "g", "", "File to get")
	flag.StringVar(&fileOut, "o", "out.dat", "Output file for retrieval")
	flag.BoolVar(&debug, "d", false, "Debug mode")
	flag.BoolVar(&check, "c", false, "Check memcached server")
	flag.Parse()

	//fmt.Printf("flags> %s, %s, %s\n", filePut, fileGet, fileOut)

	if filePut == "" && fileGet == "" && check == false {
		memza.HelpMe("Must supply file as argument (-p or -g).")
	}

	if filePut != "" {
		if err := memza.StoreFile(filePut, memcachedServer, maxFileSize, debug); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if fileGet != "" {
		if err := memza.RetrieveFile(fileGet, memcachedServer, fileOut, maxFileSize, debug); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if check == true {
		if err := memza.CheckServer(memcachedServer); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

}
