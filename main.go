package main

import (
	"fmt"
	"log"
)

func main() {

	email, err := CreateEmail()
	fmt.Println(email)
	logError(err)

	csrfToken, cookie, err := readCSRFTokenAndCookie()
	fmt.Println(csrfToken, cookie)
	logError(err)

	captcha, err := readCaptcha(cookie)
	fmt.Println(captcha)
	logError(err)

	err = register(email, captcha, csrfToken, cookie)
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
