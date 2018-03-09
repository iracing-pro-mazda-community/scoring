package score

import (
	"bufio"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iracing-pro-mazda-community/scoring/config"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/schollz/closestmatch"
)

var (
	cfg   *config.Configuration
	score map[string]int
)

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.Fatal(err)
	}

	score = make(map[string]int, 0)
}

func Print() {
	log.Printf("%#v\n", score)
}

func Match() error {
	posts, err := GetOutput()
	if err != nil {
		return err
	}

	return FuzzySearch(cfg.Tracks, posts)
}

func FuzzySearch(tracks []string, posts []forum.Post) error {
	for _, post := range posts {
		scanner := bufio.NewScanner(strings.NewReader(post.Message))
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if len(strings.Fields(text)) <= 3 {
				if len(text) > 0 {
					// try to find matches and ranking them using Levenshtein distance
					matches := fuzzy.RankFindFold(text, tracks)
					if len(matches) > 0 {
						// sort them by Levenshtein distance, pick winner
						sort.Sort(matches)
						log.Printf("DIRECTLY_MATCHED: %s: [%s] - %#v\n", post.Name, text, matches[0])
						score[matches[0].Target] = score[matches[0].Target] + 1
					} else {
						// no matches, try again with bag-of-words approach
						ClosestMatch(tracks, text, post)
					}
				}
			}
		}
	}
	return nil
}

func ClosestMatch(tracks []string, text string, post forum.Post) {
	cm := closestmatch.New(tracks, []int{2, 3})
	log.Printf("CLOSEST_MATCH: %s: [%s] - %#v\n", post.Name, text, cm.Closest(text))
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
