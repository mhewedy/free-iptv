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

func register(email string, captcha string, csrfToken string, cookie string) (err error) {

	data := url.Values{
		"token":        {csrfToken},
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

	proxyURL, err := url.Parse("http://localhost:80")
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
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

	if resp.StatusCode != http.StatusOK {
		return errors.New("http status: " + resp.Status)
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

func readCSRFTokenAndCookie() (csrfToken string, cookie string, err error) {
	resp, err := http.Get("https://my.buy-iptv.com/register.php")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	reader := resp.Body.(io.Reader)
	csrfToken, err = readCSRFToken(&reader)
	if err != nil {
		return "", "", err
	}

	cookies := strings.Join(textproto.MIMEHeader(resp.Header)["Set-Cookie"], ";")
	return csrfToken, cookies, nil
}

func readCSRFToken(r *io.Reader) (string, error) {
	document, err := goquery.NewDocumentFromReader(*r)
	if err != nil {
		return "", err
	}
	s := document.Find("#frmCheckout > input[type=hidden]:nth-child(1)")
	token, exists := s.Eq(0).Attr("value")
	if !exists {
		return "", errors.New("unable to find csrfToken")
	}
	return token, nil
}
