package forum

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/iracing-pro-mazda-community/scoring/env"
)

var cookieJar *cookiejar.Jar

func init() {
	var err error
	cookieJar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Login() error {
	location, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		panic(err)
	}
	_, utcoffset := time.Now().In(location).Zone()

	values := url.Values{}
	values.Set("username", env.MustGet("IRACING_USERNAME"))
	values.Set("password", env.MustGet("IRACING_PASSWORD"))
	values.Set("utcoffset", fmt.Sprintf("%d", utcoffset/60))
	values.Set("todaysdate", "")

	req, err := http.NewRequest("POST", "https://members.iracing.com/membersite/Login", strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Jar: cookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(string(data), "The email address or password was invalid.") ||
		resp.Header.Get("Location") == "https://members.iracing.com/membersite/failedlogin.jsp" {
		return fmt.Errorf("Login failed")
	}
	return nil
}
