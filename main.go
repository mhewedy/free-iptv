package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<MAC Address>")
		os.Exit(-1)
	}

	iptvLink, err := GetIPTVLink()
	fmt.Println(iptvLink)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdateSmartTVApp(iptvLink, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
