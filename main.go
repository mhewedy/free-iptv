package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<MAC Address>")
		os.Exit(-1)
	}

	var err error
	trial, max := 1, 20

	for {
		log.Println(">>>>> starting trial #", trial, "<<<<<")
		err = Do(os.Args[1])

		if err == nil || trial >= max {
			break
		}
		trial += 1
		duration := time.Duration(1000*trial) * time.Millisecond
		log.Println("sleeping", duration)
		time.Sleep(duration)
	}

	if err == nil {
		fmt.Println("succeed at trail #", trial)
	}
}

func Do(macAddress string) error {
	iptvLink, err := GetIPTVLink()
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(iptvLink)
	err = UpdateSmartTVApp(iptvLink, macAddress)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
