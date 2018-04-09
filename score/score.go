package score

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iracing-pro-mazda-community/scoring/config"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"github.com/schollz/closestmatch"
)

var (
	cfg   *config.Configuration
	score map[string]map[string]int64
	rx    = regexp.MustCompile(`^[[:space:] ]*([0-9A-Za-zü_\- \(\)]+)[[:space:] ]*[,  ]+[[:space:] ]*([0-4]?[0-9]+)[[:space:] ]*$`)
)

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.Fatal(err)
	}

	score = make(map[string]map[string]int64, 0)
}

func Print() {
	// csv header
	WriteToCSV(append([]string{"Driver"}, cfg.Tracks...))

	// sorted list of drivers
	var drivers []string
	for driver, _ := range score {
		drivers = append(drivers, driver)
	}
	sort.Sort(sort.StringSlice(drivers))

	ranking := make(map[string]int64, 0)
	for _, driver := range drivers {
		values := score[driver]

		// export to csv
		WriteScoreToCSV(driver, values)

		// go through tracks, if any is missing assign it maximum value to prevent skew
		for _, track := range cfg.Tracks {
			if value, ok := values[track]; ok {
				ranking[track] = ranking[track] + value
			} else {
				log.Printf("%s not voted on by %s", track, driver)
				ranking[track] = ranking[track] + int64(len(cfg.Tracks))
			}
		}
	}

	var values []int
	for _, value := range ranking {
		values = append(values, int(value))
	}
	sort.Sort(sort.IntSlice(values))

	i := 1
	for _, value := range values {
		for track, score := range ranking {
			if score == int64(value) {
				fmt.Printf("#%d: %s - [%d]\n", i, track, score)
			}
		}
		i++
	}
}

func Validate() {
	for driver, values := range score {
		//log.Printf("%s: %#v\n", driver, values)
		for i := 1; i <= 40; i++ {
			var found bool
			for _, value := range values {
				if int64(i) == value {
					if found {
						log.Printf("DUPLICATE SCORE FOR #%d by %s\n", value, driver)
						continue
					} else {
						found = true
					}
				}
			}
			if !found {
				log.Printf("NO SCORE FOR #%d FOUND by %s\n", i, driver)
			}
		}
	}
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
			vote := strings.TrimSpace(scanner.Text())
			if rx.MatchString(vote) {
				match := rx.FindStringSubmatch(vote)
				track := match[1]
				value := match[2]

				// try to find matches and ranking them using Levenshtein distance
				matches := fuzzy.RankFindFold(track, tracks)
				if len(matches) > 0 {
					// sort them by Levenshtein distance, pick winner
					sort.Sort(matches)
					//log.Printf("DIRECTLY_MATCHED: %s: [%s :: %s] - %#v\n", post.Name, track, value, matches[0])
					scoreTrack(post.Name, track, value)
				} else {
					// no matches, try again with bag-of-words approach
					ClosestMatch(tracks, track, value, post)
				}
			}
		}
	}
	return nil
}

func ClosestMatch(tracks []string, track string, value string, post forum.Post) {
	cm := closestmatch.New(tracks, []int{2, 3})
	guess := cm.Closest(track)
	log.Printf("COULD NOT SCORE: %s: [%s :: %s] - %#v\n", post.Name, track, value, guess)
	//scoreTrack(post.Name, guess, value)
}

func scoreTrack(driver, track, value string) {
	if _, ok := score[driver]; !ok {
		score[driver] = make(map[string]int64, 0)
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	score[driver][track] = score[driver][track] + v
}

func GetOutput() ([]forum.Post, error) {
	posts := make([]forum.Post, 0)

	files, err := ioutil.ReadDir("output")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := file.Name()
		if file.Name() == ".gitkeep" || filename == "result.csv" || filename == "output.txt" {
			continue
		}

		values := strings.SplitN(file.Name(), " - ", -1)
		data, err := ioutil.ReadFile(filepath.Join("output", filename))
		if err != nil {
			return nil, err
		}
		posts = append(posts, forum.Post{values[1], values[0], string(data)})
	}
	return posts, nil
}
