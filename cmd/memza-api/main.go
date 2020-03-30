package main

//
// This example is an HTTP server that can receive a file upload.
// https://www.socketloop.com/tutorials/golang-upload-file
//
// To upload a file, go to :8080/upload
// To download file, go to :8080/
//

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rcompos/memza/memza"
)

var MemcachedServer string
var maxFileSize int64 = 1024 * 1024 * 50 // 50 MB
var debug bool
var fileOut string

func main() {

	flag.StringVar(&fileOut, "o", "out.dat", "Output file for retrieval")
	flag.StringVar(&MemcachedServer, "s", os.Getenv("MEMCACHED_SERVER_URL"), "memcached_server:port")
	//flag.StringVar(&memcachedServer, "s", "localhost:11211", "memcached_server:port")
	flag.BoolVar(&debug, "d", false, "Debug mode")
	flag.Parse()

	if MemcachedServer == "" {
		MemcachedServer = "localhost:11211"
	}
	memza.MemcachedServer = MemcachedServer

	dir, err := os.Getwd()

	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	// upload -> receive
	http.HandleFunc("/upload", uploadHandler)   // Display a form for user to upload file
	http.HandleFunc("/receive", receiveHandler) // Handle incoming file

	// request -> retrieve
	http.HandleFunc("/download", downloadHandler) // Display a form for memcached key request
	http.HandleFunc("/retrieve", retrieveHandler) // memcached key retrieval

	http.HandleFunc("/test", testHandler) // Handle test
	http.HandleFunc("/info", memza.Info)  // Handle test
	http.Handle("/", http.FileServer(http.Dir(dir)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// infoHandler returns an HTML upload form
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `<html>
<head>
  <title>Memza Upload</title>
</head>
<body>
<h3>Memza - Upload a file to Memcache</h3>
<form action="/receive" method="post" enctype="multipart/form-data">
  <label for="file">Filename:</label><br>
  <input type="file" name="file" id="file"><br>
  <br>
  <input type="submit" name="submit" value="Submit">
</form>
</body>
</html>`)
	}
}

// receiveHandler accepts the file and saves it to the current working directory
func receiveHandler(w http.ResponseWriter, r *http.Request) {

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	uploadedFile := header.Filename

	//out, err := os.Create(header.Filename)
	out, err := os.Create(uploadedFile)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	//fmt.Fprintf(w, "File uploaded: ")
	//fmt.Fprintf(w, uploadedFile+"\n")

	fmt.Fprintf(w, "memcacheServer: %s\n", MemcachedServer)
	fmt.Fprintf(w, "Storing file in memcache: %s\n", uploadedFile)
	sha, errStore := memza.StoreFile(uploadedFile, MemcachedServer, maxFileSize, debug, true)
	if errStore != nil {
		fmt.Printf("error: store file in memcached failed\n")
		return
	}

	// Delete file
	errRemove := os.Remove(uploadedFile)
	if errRemove != nil {
		fmt.Printf("error: %v\n", errRemove)
	}
	fmt.Fprintf(w, "key: %s\n", uploadedFile)
	fmt.Fprintf(w, "sha256sum: %x\n", sha)

}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `<html>
<head>
  <title>Memza Download</title>
</head>
<body>
<h2>Request Memcached Key</h2>
<form action="/retrieve" method="post">
  <label for="key">memcache key:</label><br>
  <input type="text" name="key" id="key"><br>
</form>
</body>
</html>`)
	}
}

// retrieveHandler accepts the memcached key request
func retrieveHandler(w http.ResponseWriter, r *http.Request) {

	// take in POST input values
	//file, header, err := r.FormFile("file")

	requestedKey := r.FormValue("key")

	//fmt.Fprintf(w, "Retrieving file from memcache: %s\n", requestedKey)
	data, err := memza.RetrieveFile(requestedKey, MemcachedServer, fileOut, debug)
	if err != nil {
		fmt.Printf("error: retrieve file from memcached failed\n")
		return
	}
	fmt.Fprintf(w, "%s", data)

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Memza"
	w.Write([]byte(msg))
	w.Write([]byte("\n"))
}
