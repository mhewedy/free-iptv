package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/motemen/go-loghttp"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTMLElement(body io.Reader, selector string, attr string) (string, error) {
	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	s := document.Find(selector)
	el, exists := s.Eq(0).Attr(attr)
	if !exists {
		return "", errors.New("unable to find element: " + selector)
	}
	return el, nil
}

func call(method, url string, data url.Values, cookie string) (*http.Response, error) {
	c := &http.Client{Transport: &loghttp.Transport{}}

	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
