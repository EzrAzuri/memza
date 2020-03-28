package main

import (
	"flag"
	"fmt"

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

	fmt.Printf("flags> %s, %s, %s\n", filePut, fileGet, fileOut)

	if filePut == "" && fileGet == "" {
		memza.HelpMe("Must supply file as argument (-p or -g).")
	}

	if filePut != "" {
		memza.StoreFile(filePut, memcachedServer, maxFileSize)
	}

	if fileGet != "" {
		memza.RetrieveFile(fileGet, memcachedServer, fileOut, maxFileSize)
	}

	//mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	//mc := memcache.New("localhost:11211")

	//memza.CheckServer(memcachedServer)

	/*
		var zkey string = "food"
		var zval []byte = []byte("taco")
		var zfla uint32 = 77
		var zexp int32 = 4600
	*/

	//memza.Evaluator(memcachedServer, zkey, zval, zfla, zexp)

	// // memza.Setter(memcachedServer, zkey, zval, zfla, zexp)
	// // memza.Getter(memcachedServer, zkey)

}
