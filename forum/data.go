package forum

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: cookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

func GetPostsFromPage(page string, topic string) ([]Post, error) {
	if len(page) > 0 {
		topic = fmt.Sprintf("%s/%s", page, topic)
	}
	data, err := getData(fmt.Sprintf("http://members.iracing.com/jforum/posts/list/%s.page", topic))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0)
	doc.Find("tr.trPosts").Each(func(i int, s *goquery.Selection) {
		driver := s.Find("a strong").Text()
		// message, err := s.Find(".postBody").Html()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		message := s.Find(".postBody")

		id, _ := message.Attr("id")
		id = strings.TrimPrefix(id, "message")

		text := strings.Replace(message.Text(), `\n`, "\n", -1)
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\t", "", -1)
		text = strings.TrimSpace(text)
		text = strings.ToLower(text) // lowercase track names
		text = strings.Replace(text, "-", " ", -1)
		text = strings.Replace(text, "_", " ", -1)

		// go through, validate for vote and trim line by line
		var votes string
		rx := regexp.MustCompile(`^[[:space:]]*[0-9A-Za-zü_\- \(\)]+[[:space:]]*,?[[:space:]]*[0-4]?[0-9]+[[:space:]]*$`)
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			vote := strings.TrimSpace(scanner.Text())
			if rx.MatchString(vote) {
				votes += fmt.Sprintf("%s\n", vote)
			} else {
				log.Printf("Not a valid vote by [%s]: %s\n", driver, vote)
			}
		}

		posts = append(posts, Post{id, driver, votes})
	})
	return posts, nil
}

func GetAllPosts(topic string) ([]Post, error) {
	data, err := getData(fmt.Sprintf("http://members.iracing.com/jforum/posts/list/%s.page", topic))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	pages := make(map[string]string, 0)
	pages[""] = ""
	doc.Find("div.pagination a").Each(func(i int, s *goquery.Selection) {
		page, found := s.Attr("href")
		if found {
			page = strings.TrimPrefix(page, "/jforum/posts/list/")
			page = strings.TrimSuffix(page, fmt.Sprintf("/%s.page", topic))
			pages[page] = page
		}
	})

	posts := make([]Post, 0)
	for page := range pages {
		ps, err := GetPostsFromPage(page, topic)
		if err != nil {
			return nil, err
		}
		posts = append(posts, ps...)
	}
	return posts, nil
}
