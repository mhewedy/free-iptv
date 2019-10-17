package main

import (
	"fmt"
	"log"
)

func main() {

	iptvLink, err := GetIPTVLink()
	fmt.Println(iptvLink)
	if err != nil {
		log.Fatal(err)
	}
}
