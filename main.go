package main

import (
	"fmt"
	"log"
)

func main() {

	iptvLink, err := GetIPTVLink()
	fmt.Println(iptvLink)
	logError(err)
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
