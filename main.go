package main

import (
	"log"

	"github.com/iracing-pro-mazda-community/scoring/config"
	"github.com/iracing-pro-mazda-community/scoring/forum"
	"github.com/iracing-pro-mazda-community/scoring/score"
)

func main() {
	// download/update votes
	if config.Get().Download {
		if err := forum.Login(); err != nil {
			log.Fatal(err)
		}

		log.Println("Download voting data ...")
		posts, err := forum.GetAllPosts(config.Get().Topic)
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
	if config.Get().Score {
		log.Println("Score votes ...")
		if err := score.Match(); err != nil {
			log.Fatal(err)
		}
		score.Validate()
		score.Print()
	}
}
