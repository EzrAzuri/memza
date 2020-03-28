package memza

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"

	//"io"
	"io/ioutil"
	"log"
	"os"

	//"strings"
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

// RetrieveFile get file contents for given key filename
func RetrieveFile(file, mserver, outFile string, max int64) error {
	// Get number of required chunks for file
	num, shasum, err := numChunks(file, max)
	if err != nil {
		fmt.Println("ERROR: %v\n", err)
	}
	fmt.Printf("chunks: %v\n", num)
	fmt.Printf("sha256sum:\n %x\n", shasum)

	// Single-chunk
	if num == 1 {
		singleChunk, err := Getter(mserver, file)
		check(err)

		errWrite := ioutil.WriteFile(outFile, singleChunk, 0644)
		check(errWrite)
	} else if num > 1 {

		// Multi-chunk
		// reconstitute
		// write file

	}

	// Check file sums
	data, err := ioutil.ReadFile(file)
	check(err)

	// Hash the file and output results
	fmt.Printf("SHA-256: %x\n", sha256.Sum256(data))
	return err

}

// StoreFile key: filename, value: file contents
func StoreFile(file, mserver string, max int64) error {

	// Get number of required chunks for file
	num, shasum, err := numChunks(file, max)
	if err != nil {
		fmt.Println("ERROR: %v\n", err)
	}
	fmt.Printf("chunks: %v\n", num)
	fmt.Printf("sha256sum: %x\n", shasum)

	if num == 1 { // Single-chunk
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		// Set key with file bytes as value
		// if file fits, it sits - small file under 1MB
		errSetter := Setter(mserver, file, fileBytes, 0, 0)
		if errSetter != nil {
			log.Fatal(err)
		}
	} else if num > 1 {

		fmt.Println("--> Needs chunking!")
		// open file
		// defer close
		for i := 1; i <= num; i++ {

			fmt.Printf("\tChunk: %s\n", i)
			//
			// read 1MB
			// set mcache value
			//
		}

	}
	return err

}

// numChunks determine number of chunks needed
func numChunks(fileName string, max int64) (int, []byte, error) {

	//mc := memcache.New(mserver)
	// Set key
	//fmt.Printf("Evaluate key -> %s\tvalue: %s\tflag: %d\texp: %d\n", key, val, fla, exp)
	fmt.Printf("File: %s\n", fileName)
	sizeBytes := fileSize(fileName)
	fmt.Printf("Size (bytes): %d\n", sizeBytes)

	// Empty file check
	if sizeBytes == 0 {
		return 0, []byte{}, errors.New("Zero file size!")
	}

	// Max size check - 50MB
	if sizeBytes > max {
		fmt.Printf("Max size: %d\n", max)
		errMsg := fmt.Sprintf("ERROR: File too large: %d\n", sizeBytes)
		return 0, []byte{}, errors.New(errMsg)
	}

	/*
		input := strings.NewReader(fileName)
		hash := sha256.New()
		if _, err := io.Copy(hash, input); err != nil {
			log.Fatal(err)
		}
		sum := hash.Sum(nil)
	*/

	data, err := ioutil.ReadFile(fileName)
	check(err)

	fileSHA256 := sha256.Sum256(data)

	// Hash the file and output results
	//fmt.Printf("SHA-256: %x\n", sha256.Sum256(data))
	fmt.Printf("SHA-256: %x\n", fileSHA256)

	// Figure out how many 1MB chunks
	floatChunks := float64(sizeBytes) / (1024 * 1024)
	fmt.Printf("Chunks (1MB) Float: %f\n", floatChunks)

	//intChunks := math.Floor(floatChunks)
	intChunks := int(floatChunks)
	fmt.Printf("Chunks (1MB) Int: %d\n", intChunks)

	if floatChunks > float64(intChunks) {
		intChunks += 1
	}

	return intChunks, fileSHA256[:], err

}

// Setter is good for setting mcache values
func Setter(mserver, key string, val []byte, fla uint32, exp int32) error {
	mc := memcache.New(mserver)
	// Set key
	//fmt.Printf("Set key -> %s\tvalue: %s\tflag: %d\texp: %d\n", key, val, fla, exp)
	//mc.Set(&memcache.Item{Key: "foo", Value: []byte("baarr")})
	//mc.Set(&memcache.Item{Key: key, Value: []byte(val), Flags: fla, Expiration: exp})
	err := mc.Set(&memcache.Item{Key: key, Value: val, Flags: fla, Expiration: exp})
	return err
}

// Getter is good for getting mcache values
func Getter(mserver, key string) ([]byte, error) {
	mc := memcache.New(mserver)
	// Get key
	fmt.Printf("Get key -> %s\n", key)
	myitem, err := mc.Get(key)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	/*
		fmt.Printf("item: %v\n", myitem)
		fmt.Printf("key: %v\n", myitem.Key)
		fmt.Printf("value: %v\n", string(myitem.Value))
		fmt.Printf("flags: %d\n", myitem.Flags)
		//fmt.Printf("expiration: %d\n", myitem.Expiration)
	*/
	return myitem.Value, err
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

// HelpMe provides help usage message
func HelpMe(msg string) {
	if msg != "" {
		fmt.Printf("%s\n\n", msg)
	}
	fmt.Println("Supply file name i.e. /path/to/myfile.txt")
	flag.PrintDefaults()
	os.Exit(1)
}

// fileSize checks file size
func fileSize(f string) int64 {
	fi, err := os.Stat(f)
	if err != nil {
		log.Fatal(err)
	}
	return fi.Size()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
