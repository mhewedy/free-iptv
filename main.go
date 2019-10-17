package main

import (
	"fmt"
	"log"
)

func main() {

	email, err := CreateEmail()
	fmt.Println(email)
	logError(err)

	token, cookie, err := readTokenAndCookie()
	fmt.Println(token, cookie)
	logError(err)

	captcha, err := readCaptcha(cookie)
	fmt.Println(captcha)
	logError(err)

	err = register(email, captcha, token, cookie)
	logError(err)
}

func CreateEmail() (email string, err error) {
	return ReadDOMElement("https://10minutemail.net/", "#fe_text", "value")
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
