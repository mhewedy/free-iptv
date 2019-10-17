package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/motemen/go-loghttp"
	"github.com/otiai10/gosseract"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

func GetIPTVLink() (string, error) {

	email, err := CreateEmail()
	if err != nil {
		return "", err
	}

	password := "p@ssw0rd"
	cookie, err := DoRegister(email, password)
	if err != nil {
		return "", err
	}

	err = Buy(email, password, cookie)
	if err != nil {
		return "", err
	}

	link, err := GetM3ULink(cookie)
	if err != nil {
		return "", err
	}

	return link, nil
}

func CreateEmail() (email string, err error) {
	return ReadDOMElement("https://10minutemail.net/", "#fe_text", "value")
}

func DoRegister(email string, password string) (cookie string, err error) {
	log.Println("start registration using email: ", email)

	token, cookie, err := readTokenAndCookie("https://my.buy-iptv.com/register.php",
		"#frmCheckout > input[type=hidden]:nth-child(1)", "")
	if err != nil {
		return "", err
	}

	captcha, err := readCaptcha(cookie)
	if err != nil {
		return "", err
	}
	log.Println("captcha value: ", captcha)

	err = register(email, password, captcha, token, cookie)
	if err != nil {
		return "", err
	}

	fmt.Println("cookie:", cookie)

	return cookie, nil
}

func register(email string, password string, captcha string, token string, cookie string) (err error) {

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

	resp, err := call("POST", "https://my.buy-iptv.com/register.php", data, cookie)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path != "/clientarea.php" {
		bytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bytes))
		return errors.New("unable to create account")
	}
	return nil
}

func Buy(email string, password string, cookie string) error {

	log.Println("start buying 1-day token")

	// Add to cart
	_, err := call("POST", "https://my.buy-iptv.com/cart.php?a=add&pid=88", nil, cookie)
	if err != nil {
		return err
	}

	// Buy the cart items
	token, _, err := readTokenAndCookie("https://my.buy-iptv.com/cart.php?a=view",
		"#order-cartx > div.accout-row > div.col-md-5.total-bar > form > input[type=hidden]", cookie)
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

	resp, err := call("POST", "https://my.buy-iptv.com/cart.php?a=checkout", data, cookie)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path != "/cart.php" || resp.Request.URL.RawQuery != "a=complete" {
		bytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bytes))
		return errors.New("unable to buy")
	}
	return nil
}

func GetM3ULink(cookie string) (string, error) {

	// Get manage link
	resp, err := call("GET", "https://my.buy-iptv.com/clientarea.php", nil, cookie)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	s := document.Find("#products > div:nth-child(1) > div > div.activ-right > ul > li.manag > a")
	manageLink, exists := s.Eq(0).Attr("href")
	if !exists {
		return "", errors.New("unable to find box link")
	}
	manageLink = "https://my.buy-iptv.com" + manageLink

	// Get details link
	token, _, err := readTokenAndCookie(manageLink,
		"#tabChangepw > div > div > div.files-body > form > input[type=hidden]:nth-child(1)", cookie)
	if err != nil {
		return "", err
	}

	parse, err := url.Parse(manageLink)
	if err != nil {
		return "", err
	}

	// Get m3u link
	resp, err = call("POST", "https://my.buy-iptv.com/clientarea.php?action=productdetails", url.Values{
		"token":        {token},
		"customAction": {"manage"},
		"id":           {parse.Query().Get("id")},
	}, cookie)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	document, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	s = document.Find("#m3ulinks")
	m3uLink, exists := s.Eq(0).Attr("value")
	if !exists {
		return "", errors.New("unable to find m3u link")
	}

	return m3uLink, nil
}

func readTokenAndCookie(pageURL string, tokenSelector string, inCookie string) (token string, cookie string, err error) {

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

	resp, err := call("GET", pageURL, nil, inCookie)
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

	return string(bytes), nil
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
