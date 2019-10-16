package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func main() {
	email, err := CreateEmail()
	logError(err)

	fmt.Println(email)
}

func CreateEmail() (email string, err error) {
	resp, err := http.Get("https://10minutemail.net/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	s := document.Find("#fe_text")
	if email, exists := s.Eq(0).Attr("value"); exists {
		return email, nil
	} else {
		return "", errors.New("cannot create email")
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
