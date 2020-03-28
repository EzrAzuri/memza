package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rcompos/memza/memza"
)

func main() {

	var memcachedServer string = "localhost:11211"
	var filePut, fileGet, fileOut string
	flag.StringVar(&filePut, "p", "", "File to put")
	flag.StringVar(&fileGet, "g", "", "File to get")
	flag.StringVar(&fileOut, "o", "out.dat", "Output file for retrieval")
	flag.Parse()
	var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB

	//fmt.Printf("flags> %s, %s, %s\n", filePut, fileGet, fileOut)

	if filePut == "" && fileGet == "" {
		memza.HelpMe("Must supply file as argument (-p or -g).")
	}

	if filePut != "" {
		if err := memza.StoreFile(filePut, memcachedServer, maxFileSize); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if fileGet != "" {
		if err := memza.RetrieveFile(fileGet, memcachedServer, fileOut, maxFileSize); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

}
