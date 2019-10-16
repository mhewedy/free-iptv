package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/otiai10/gosseract"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	//email, err := CreateEmail()
	//logError(err)

	CreateIPTVAccount("")
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

func CreateIPTVAccount(email string) {

	readCaptcha := func() ([]byte, error) {
		resp, err := http.Get("https://my.buy-iptv.com/includes/verifyimage.php")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}

	bytes, _ := readCaptcha()
	_ = ioutil.WriteFile(`abc.png`, bytes, os.ModePerm)

	client := gosseract.NewClient()
	defer client.Close()
	_ = client.SetImageFromBytes(bytes)
	text, _ := client.Text()
	fmt.Println(text)

}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
