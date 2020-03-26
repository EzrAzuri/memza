package main

import(
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {

	fmt.Println("Zapier Interview")

	//mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	mc := memcache.New("localhost:11211")

	// Check connection to memcached server
	errPing := mc.Ping()
	if errPing != nil {
		fmt.Println("ERROR: %v", errPing)
	}
	fmt.Printf("ping successfull!\n")

	// Set key
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("baarr")})

	// Get key
	myitem, err := mc.Get("foo")
	if err != nil {
		fmt.Println("ERROR: %v", err)
	}

	fmt.Printf("item: %v\n", myitem)
	fmt.Printf("key: %v\n", myitem.Key)
	fmt.Printf("value: %v\n", string(myitem.Value))
	fmt.Printf("flags: %v\n", myitem.Flags)
	fmt.Printf("expiration: %v\n", myitem.Expiration)

}