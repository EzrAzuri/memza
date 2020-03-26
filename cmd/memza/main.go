package main

import (
	"fmt"

	"github.com/rcompos/memza/memza"
)

var memcachedServer string = "localhost:11211"

func main() {

	fmt.Println("Memza")

	//mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	//mc := memcache.New("localhost:11211")

	memza.CheckServer(memcachedServer)

	fmt.Printf("azmeM\n")

}
