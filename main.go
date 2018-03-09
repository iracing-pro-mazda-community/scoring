package main

import (
	"bufio"
	"log"
	"strings"

	"github.com/iracing-pro-mazda-community/scoring/env"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/schollz/closestmatch"
)

func main() {
	if err := forum.Login(); err != nil {
		log.Fatal(err)
	}

	posts, err := forum.GetAllPosts(env.MustGet("IRACING_TOPIC"))
	if err != nil {
		log.Fatal(err)
	}

	tracks := []string{
		"Laguna Seca",
		"Okayama",
		"Summit Point",
		"Lime Rock",
		"Daytona 2007",
		"Road America",
		"Spa",
		"Road Atlanta",
		"Watkins",
		"Monza",
		"Interlagos",
		"Suzuka",
		"Mosport",
		"Montreal",
		"Donington",
		"Nordschleife",
		"Bathurst",
		"Nurbürgring GP",
		"Imola",
		"Le Mans",
		"Mid Ohio",
		"Phillip Island",
		"Silverstone",
		"Oulton Park",
		"Brands Hatch",
		"Sebring",
		"Barber",
		"Zolder",
		"Zandvoort",
		"Snetterton",
		"Indianapolis",
		"Sonoma",
		"Miami",
		"COTA",
		"Motegi",
		"VIR",
		"Oran Park",
	}
	cm := closestmatch.New(tracks, []int{2, 3})

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
}
