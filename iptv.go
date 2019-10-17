package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/otiai10/gosseract"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

const password = "p@ssw0rd"

func readTokenAndCookie() (token string, cookie string, err error) {

	readToken := func(r io.Reader) (string, error) {
		document, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			return "", err
		}
		s := document.Find("#frmCheckout > input[type=hidden]:nth-child(1)")
		token, exists := s.Eq(0).Attr("value")
		if !exists {
			return "", errors.New("unable to find token")
		}
		return token, nil
	}

	resp, err := http.Get("https://my.buy-iptv.com/register.php")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	token, err = readToken(resp.Body.(io.Reader))
	if err != nil {
		return "", "", err
	}

	cookies := strings.Join(textproto.MIMEHeader(resp.Header)["Set-Cookie"], ";")
	return token, cookies, nil
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
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", "https://my.buy-iptv.com/register.php",
		strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return errors.New("unable to create account: " +
			"http status should equal 302, however it is : " + resp.Status)
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
