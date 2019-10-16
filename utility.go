package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func ReadDOMElement(url string, selector string, attr string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	s := document.Find(selector)
	if email, exists := s.Eq(0).Attr(attr); exists {
		return email, nil
	} else {
		return "", errors.New("unable to find attr")
	}
}
