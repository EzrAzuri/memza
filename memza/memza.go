package memza

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	//"math"

	"github.com/bradfitz/gomemcache/memcache"
)

//var memcachedServer string = "localhost:11211"

// Chunker splits the file into parts
func Chunker(f string) {
	fmt.Println(f)
}

// UnChunker recombines the parts into file
func UnChunker(f string) {
	fmt.Println(f)
}

// fileSize checks file size
func fileSize(f string) int64 {
	fi, err := os.Stat(f)
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}

// NumChunks is good for evaluating mcache values
func NumChunks(fileName string, max int64) (int, error) {

	//mc := memcache.New(mserver)
	// Set key
	//fmt.Printf("Evaluate key -> %s\tvalue: %s\tflag: %d\texp: %d\n", key, val, fla, exp)
	fmt.Printf("File: %s\n", fileName)
	sizeBytes := fileSize(fileName)
	fmt.Printf("Size (bytes): %d\n", sizeBytes)

	// Empty file check
	if sizeBytes == 0 {
		return 0, errors.New("Zero file size!")
	}

	// Max size check - 50MB
	if sizeBytes > max {
		fmt.Printf("Max size: %d\n", max)
		errMsg := fmt.Sprintf("ERROR: File too large: %d\n", sizeBytes)
		return 0, errors.New(errMsg)
	}

	input := strings.NewReader(fileName)
	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)

	fmt.Printf("SHA-256: %x\n", sum)

	// Figure out how many 1MB chunks
	floatChunks := float64(sizeBytes) / (1024 * 1024)
	fmt.Printf("Chunks (1MB) Float: %f\n", floatChunks)

	//intChunks := math.Floor(floatChunks)
	intChunks := int(floatChunks)
	fmt.Printf("Chunks (1MB) Int: %d\n", intChunks)

	if floatChunks > float64(intChunks) {
		intChunks += 1
	}

	return intChunks, nil

}

// Setter is good for setting mcache values
func Setter(mserver, key string, val []byte, fla uint32, exp int32) {
	mc := memcache.New(mserver)
	// Set key
	fmt.Printf("Set key -> %s\tvalue: %s\tflag: %d\texp: %d\n", key, val, fla, exp)
	//mc.Set(&memcache.Item{Key: "foo", Value: []byte("baarr")})
	//mc.Set(&memcache.Item{Key: key, Value: []byte(val), Flags: fla, Expiration: exp})
	mc.Set(&memcache.Item{Key: key, Value: val, Flags: fla, Expiration: exp})
}

// Getter is good for getting mcache values
func Getter(mserver, key string) {
	mc := memcache.New(mserver)
	// Get key
	fmt.Printf("Get key -> %s\n", key)
	myitem, err := mc.Get(key)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	fmt.Printf("item: %v\n", myitem)
	fmt.Printf("key: %v\n", myitem.Key)
	fmt.Printf("value: %v\n", string(myitem.Value))
	fmt.Printf("flags: %d\n", myitem.Flags)
	//fmt.Printf("expiration: %d\n", myitem.Expiration)

}

// CheckServer memcached server status
func CheckServer(memcachedServer string) {

	fmt.Println("Memza->CheckServer->")

	//mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	//mc := memcache.New("localhost:11211")
	mc := memcache.New(memcachedServer)

	// Check connection to memcached server
	fmt.Printf("Ping memcached server\n")
	errPing := mc.Ping()
	if errPing != nil {
		fmt.Printf("ping failed!\n")
		fmt.Println("ERROR: %v", errPing)
	}
	fmt.Printf("ping successfull!\n")

	// Set key
	keyIn := "foo"
	valIn := "baarrr"
	fmt.Printf("Set Item\n")
	fmt.Printf("Set key -> %s\tvalue: %s\n", keyIn, valIn)
	//mc.Set(&memcache.Item{Key: "foo", Value: []byte("baarr")})
	mc.Set(&memcache.Item{Key: keyIn, Value: []byte(valIn)})

	// Get key
	fmt.Printf("Get key ->\n")
	myitem, err := mc.Get("foo")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	fmt.Printf("item: %v\n", myitem)
	fmt.Printf("key: %v\n", myitem.Key)
	fmt.Printf("value: %v\n", string(myitem.Value))
	fmt.Printf("flags: %v\n", myitem.Flags)
	fmt.Printf("expiration: %v\n", myitem.Expiration)

}

func HelpMe(msg string) {
	if msg != "" {
		fmt.Printf("%s\n\n", msg)
	}
	fmt.Println("Supply file name i.e. /path/to/myfile.txt")
	flag.PrintDefaults()
	os.Exit(1)
}
