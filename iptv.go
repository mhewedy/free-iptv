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
	"os"
	"strings"
)

const password = "p@ssw0rd"

func CreateEmail() (email string, err error) {
	return ReadDOMElement("https://10minutemail.net/", "#fe_text", "value")
}

func GetIPTVLink() (string, error) {

	email, cookie, err := DoRegister()
	if err != nil {
		return "", err
	}

	err = Buy(email, password, cookie)
	if err != nil {
		return "", err
	}

	return "<TODO>", nil
}

func DoRegister() (email string, cookie string, err error) {
	log.Println("start registration")

	email, err = CreateEmail()
	if err != nil {
		return "", "", err
	}
	log.Println("register using using email: ", email)

	token, cookie, err := readTokenAndCookie("https://my.buy-iptv.com/register.php")
	if err != nil {
		return "", "", err
	}

	captcha, err := readCaptcha(cookie)
	if err != nil {
		return "", "", err
	}
	log.Println("captcha value: ", captcha)

	err = register(email, captcha, token, cookie)
	if err != nil {
		return "", "", err
	}

	fmt.Println("cookie:", cookie)

	return email, cookie, nil
}

func readTokenAndCookie(pageURL string) (token string, cookie string, err error) {

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
		_ = resp.Header.Write(os.Stdout)
		return errors.New("unable to create account")
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

func Buy(email string, password string, cookie string) error {

	log.Println("start buying 1-day token")

	token, _, err := readTokenAndCookie("https://my.buy-iptv.com/cart.php?a=view")
	if err != nil {
		return err
	}

	data := url.Values{
		"token":         {token},
		"submit":        {"true"},
		"custtype":      {"existing"},
		"loginemail":    {email},
		"loginpassword": {password},
		"firstname":     {"Abbosa"},
		"lastname":      {"fornasa"},
		"email":         {email},
		"phonenumber":   {"7185511111"},
		"address1":      {"1th Avenue"},
		"city":          {"Springfield Gardens"},
		"state":         {"Alabama"},
		"postcode":      {"12345"},
		"paymentmethod": {"banktransfer"},
		"ccinfo":        {"new"},
		"cctype":        {"visa"},
		"ccnumber":      {""},
		"ccexpirydate":  {""},
		"cccvv":         {""},
		"cccvvexisting": {""},
		"accepttos":     {"on"},
		"notes":         {""},
	}
	c := &http.Client{}
	req, err := http.NewRequest("POST", "https://my.buy-iptv.com/cart.php?a=checkout",
		strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path != "/cart.php" && resp.Request.URL.RawQuery != "a=complete" {
		return errors.New("unable to buy: " +
			"the request redirected to  : " + resp.Request.RequestURI)
	}
	return nil
}
