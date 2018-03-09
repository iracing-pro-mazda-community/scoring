package score

import (
	"bufio"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/iracing-pro-mazda-community/scoring/config"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/schollz/closestmatch"
)

var cfg *config.Configuration

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.Fatal(err)
	}
}

func Match() error {
	posts, err := GetOutput()
	if err != nil {
		return err
	}

	cm := closestmatch.New(cfg.Tracks, []int{2, 3})
	for _, post := range posts {
		scanner := bufio.NewScanner(strings.NewReader(post.Message))
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if len(strings.Fields(text)) <= 3 {
				if len(text) > 0 {
					log.Printf("%s: [%s] - %#v\n", post.Name, text, cm.Closest(text))
				}
			}
		}
	}
	return nil
}

func GetOutput() ([]forum.Post, error) {
	posts := make([]forum.Post, 0)

	files, err := ioutil.ReadDir("output")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Name() == ".gitkeep" {
			continue
		}

		values := strings.SplitN(file.Name(), " - ", -1)
		data, err := ioutil.ReadFile(filepath.Join("output", file.Name()))
		if err != nil {
			return nil, err
		}
		posts = append(posts, forum.Post{values[1], values[0], string(data)})
	}
	return posts, nil
}
