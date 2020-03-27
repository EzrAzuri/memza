package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rcompos/memza/memza"
)

var memcachedServer string = "localhost:11211"

func main() {

	var fileName string
	flag.StringVar(&fileName, "f", os.Getenv("MEMZA_FILE"), "File to store")
	flag.Parse()

	if fileName == "" {
		memza.HelpMe("Must supply file as argument (-f).")
	}

	var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB

	// Get number of required chunks for file
	num, err := memza.NumChunks(fileName, maxFileSize)
	if err != nil {
		fmt.Println("ERROR: %v\n", err)
	}
	fmt.Printf("chunks: %v\n", num)

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
