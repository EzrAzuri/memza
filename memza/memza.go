package memza

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	//"strings"
	"io"
	//"math"

	"github.com/bradfitz/gomemcache/memcache"
)

const devilsBytes int = 61
const bufferSizeMax int = 1024*1024 - devilsBytes

// RetrieveFile get file contents for given key filename
func RetrieveFile(f, mserver, outFile string, max int64, dbug bool) error {

	// Get number of required chunks for file
	filehash, num, errGet := Getter(mserver, f, dbug)
	check(errGet)

	if dbug == true {
		fmt.Printf("RetrieveFile ->\n")
		fmt.Printf("Key: %s\n", f)
		fmt.Printf("Chunks: %v\n", num)
		fmt.Printf("Filehash: %x\n", string(filehash))
	}

	// open file
	file, errCreate := os.Create(outFile)
	check(errCreate)
	defer file.Close()

	// reconstitute
	for i := 1; i <= int(num); i++ {
		chunkKey := f + "-" + strconv.Itoa(i)
		// Get single chunk
		chunkItem, _, err := Getter(mserver, chunkKey, dbug)
		check(err)
		// write file
		n2, werr := file.Write(chunkItem)
		if dbug == true {
			fmt.Printf("chunkKey: %s\n", chunkKey)
			fmt.Printf("\tchunk: %v\n", i)
			fmt.Printf("\tBytes written: %d\n", n2)
			check(werr)
		}
	}

	// Read newly created file
	data, errRead := ioutil.ReadFile(outFile)
	check(errRead)
	// Hash the file and output results
	newHash := sha256.Sum256(data)

	fmt.Printf("%s %x\n", outFile, newHash)

	//badHash := []byte{'1', '9', 'a', 'f'} // For TESTING ONLY
	//compareResult := bytes.Compare(filehash[:], badHash)
	compareResult := bytes.Compare(filehash[:], newHash[:])
	//fmt.Printf("Hash compare: %v\n", compareResult)
	var err error
	if compareResult != 0 {
		err = errors.New("hash mismatch")
	}

	return err

}

// StoreFile key: filename, value: file contents
func StoreFile(f, mserver string, max int64, dbug bool) error {

	bufferSize := bufferSizeMax - len(f)

	// Get number of required chunks for file
	num, shasum, err := numChunks(f, bufferSize, max, dbug)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	if dbug == true {
		fmt.Printf("StoreFile->\n")
		fmt.Printf("\tKey: %s\n", f)
		fmt.Printf("\tValue: %x\n", shasum)
		fmt.Printf("chunks: %v\n", num)
		fmt.Printf("sha256sum: %x\n", shasum)
		fmt.Printf("Setting item:\n")
	}

	// Set key named after file with value of shasum
	errSetterFile := Setter(mserver, f, shasum[:], uint32(num), 0)
	check(errSetterFile)

	// open file
	file, err := os.Open(f)
	check(err)
	defer file.Close()

	buffer := make([]byte, bufferSize)
	var i int = 1
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		buff := buffer[:bytesread]

		// Set file contents
		//fileKey := f + "-" + strconv.Itoa(i) + "-" + hex.EncodeToString(shasum)
		fileKey := f + "-" + strconv.Itoa(i)

		if dbug == true {
			fmt.Printf("\tChunk: %v\n", i)
			fmt.Println("\tBytes read: ", bytesread)
			fmt.Printf("\tKey: %v\n", fileKey)
		}

		errSetterHash := Setter(mserver, fileKey, buff, 0, 0)
		check(errSetterHash)

		i++
	}

	return err

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
func Getter(mserver, key string, dbug bool) ([]byte, uint32, error) {
	mc := memcache.New(mserver)
	// Get key
	if dbug == true {
		fmt.Printf("Get key -> %s\n", key)
	}
	myitem, err := mc.Get(key)
	if err != nil {
		fmt.Printf("%v", err)
	}

	/*
		fmt.Printf("item: %v\n", myitem)
		fmt.Printf("key: %v\n", myitem.Key)
		fmt.Printf("value: %v\n", string(myitem.Value))
		fmt.Printf("flags: %d\n", myitem.Flags)
		//fmt.Printf("expiration: %d\n", myitem.Expiration)
	*/

	return myitem.Value, myitem.Flags, err
}

// CheckServer memcached server status
func CheckServer(memcachedServer string) error {

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

	return err

}

// numChunks determine number of chunks needed
func numChunks(fileName string, chunksize int, max int64, dbug bool) (int, [32]byte, error) {

	//mc := memcache.New(mserver)
	// Set key
	//fmt.Printf("Evaluate key -> %s\tvalue: %s\tflag: %d\texp: %d\n", key, val, fla, exp)
	sizeBytes := fileSize(fileName)

	// Empty file check
	if sizeBytes == 0 {
		return 0, [32]byte{}, errors.New("Zero file size!")
	}

	// Max size check - 50MB
	if sizeBytes > max {
		fmt.Printf("Max size: %d\n", max)
		errMsg := fmt.Sprintf("ERROR: File too large: %d\n", sizeBytes)
		return 0, [32]byte{}, errors.New(errMsg)
	}

	data, err := ioutil.ReadFile(fileName)
	check(err)

	fileSHA256 := sha256.Sum256(data)

	// Hash the file and output results
	//fmt.Printf("SHA-256: %x\n", sha256.Sum256(data))

	// Figure out how many 1MB chunks
	//floatChunks := float64(sizeBytes) / (1024 * 1024)
	floatChunks := float64(sizeBytes) / float64(chunksize)

	//intChunks := math.Floor(floatChunks)
	intChunks := int(floatChunks)
	if floatChunks > float64(intChunks) {
		intChunks += 1
	}

	if dbug == true {
		fmt.Printf("File: %s\n", fileName)
		fmt.Printf("Size (bytes): %d\n", sizeBytes)
		fmt.Printf("SHA-256: %x\n", fileSHA256)
		fmt.Printf("Chunks (1MB) Float: %f\n", floatChunks)
		fmt.Printf("Chunks (1MB) Int: %d\n", intChunks)
	}

	//return intChunks, fileSHA256[:], err
	return intChunks, fileSHA256, err

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
