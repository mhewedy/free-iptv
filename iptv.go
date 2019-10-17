package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/otiai10/gosseract"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

const password = "p@ssw0rd"

func CreateEmail() (email string, err error) {
	return ReadDOMElement("https://10minutemail.net/", "#fe_text", "value")
}

func GetIPTVLink() (string, error) {

	_, err := DoRegister()
	if err != nil {
		return "", err
	}

	/*cookie, err := DoLogin(email, password)
	if err != nil {
		return "", err
	}

	fmt.Println("Login cookie:", cookie)*/

	return "<TODO>", nil
}

func DoRegister() (string, error) {
	log.Println("start registration")

	email, err := CreateEmail()
	if err != nil {
		return "", err
	}
	log.Println("register using using email: ", email)
	token, cookie, err := readTokenAndCookie("https://my.buy-iptv.com/register.php",
		"#frmCheckout > input[type=hidden]:nth-child(1)")
	if err != nil {
		return "", err
	}
	captcha, err := readCaptcha(cookie)
	if err != nil {
		return "", err
	}
	log.Println("captcha read as: ", captcha)

	err = register(email, captcha, token, cookie)
	if err != nil {
		return "", err
	}

	fmt.Println("cookie:", cookie)

	return email, nil
}

func readTokenAndCookie(pageURL string, tokenSelector string) (token string, cookie string, err error) {

	readToken := func(r io.Reader) (string, error) {
		document, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			return "", err
		}
		s := document.Find(tokenSelector)
		token, exists := s.Eq(0).Attr("value")
		if !exists {
			return "", errors.New("unable to find token")
		}
		return token, nil
	}

	resp, err := http.Get(pageURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	token, err = readToken(resp.Body.(io.Reader))
	if err != nil {
		return "", "", err
	}

	cookie = strings.Join(textproto.MIMEHeader(resp.Header)["Set-Cookie"], ";")
	return token, cookie, nil
}

func register(email string, captcha string, token string, cookie string) (err error) {

	data := url.Values{
		"token":        {token},
		"register":     {"true"},
		"firstname":    {"Abbosa"},
		"lastname":     {"fornasa"},
		"email":        {email},
		"address1":     {"1th Avenue"},
		"phonenumber":  {"7185511111"},
		"city":         {"Springfield Gardens"},
		"state":        {"Alabama"},
		"postcode":     {"12345"},
		"country":      {"US"},
		"address2":     {},
		"password":     {password},
		"password2":    {password},
		"securityqid":  {"1"},
		"securityqans": {"1"},
		"code":         {captcha},
		"accepttos":    {"on"},
		"companyname":  {},
	}
	c := &http.Client{}
	req, err := http.NewRequest("POST", "https://my.buy-iptv.com/register.php",
		strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path != "/clientarea.php" {
		return errors.New("unable to create account: " +
			"the request redirected to  : " + resp.Request.RequestURI)
	}

	return nil
}

func readCaptcha(cookie string) (captcha string, err error) {

	c := &http.Client{}
	req, err := http.NewRequest("GET", "https://my.buy-iptv.com/includes/verifyimage.php", nil)

	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	client := gosseract.NewClient()
	defer client.Close()
	_ = client.SetImageFromBytes(bytes)
	text, _ := client.Text()
	return text[:5], nil
}
