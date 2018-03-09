package main

import (
	"log"

	"github.com/iracing-pro-mazda-community/scoring/config"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/iracing-pro-mazda-community/scoring/score"
)

var cfg *config.Configuration

func init() {
	var err error
	cfg, err = config.Get()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := forum.Login(); err != nil {
		log.Fatal(err)
	}

	// download/update votes
	if cfg.Download {
		log.Println("Download voting data ...")
		posts, err := forum.GetAllPosts(cfg.Topic)
		if err != nil {
			log.Fatal(err)
		}

		// save them all
		for _, post := range posts {
			if err := post.Store(); err != nil {
				log.Fatal(err)
			}
		}
	}

	// match votes
	log.Println("Score votes ...")
	if err := score.Match(); err != nil {
		log.Fatal(err)
	}
	score.Print()
}
