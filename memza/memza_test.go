package memza

// Test

import (
	"encoding/hex"
	"fmt"
	"testing"
)

/*
Expected results here:
*/

func TestStoreFile(t *testing.T) {

	file := "CyberOwl.jpeg"
	memcachedServer := "0.0.0.0:11211"
	var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB
	var debug bool = true
	var force bool = true

	sha256sumString := "e48aeb79e0d0e7f50d37f9c8d1af527c4cd87f0749e01925e10b9f2e2cd14cfa"
	sha256sum, errHex := hex.DecodeString(sha256sumString)
	if errHex != nil {
		fmt.Printf("error: hex decode string failed\n")
		fmt.Printf("%v\n", errHex)
		return
	}

	fmt.Printf("Running StoreFile\n")
	sha, err := StoreFile(file, memcachedServer, maxFileSize, debug, force)
	if err != nil {
		fmt.Printf("error: store file in memcached failed\n")
		return
	}
	fmt.Printf("Checking file sums\n")
	if string(sha[:]) != string(sha256sum) {
		fmt.Printf("error: sha256sums do not match\n")
		return
	}

}
