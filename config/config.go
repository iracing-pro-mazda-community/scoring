package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	cfg  *Configuration = nil
	once sync.Once
)

type Configuration struct {
	Topic         string
	Tracks        []string
	Download      bool
	Score         bool
	OverwriteVote bool `json:"overwrite_vote"`
}

func Get() *Configuration {
	once.Do(func() {
		cfg = &Configuration{}

		file, err := os.Open("configuration.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			log.Fatal(err)
		}

		// lowercase track names
		for t := range cfg.Tracks {
			cfg.Tracks[t] = strings.ToLower(cfg.Tracks[t])
			cfg.Tracks[t] = strings.Replace(cfg.Tracks[t], "-", " ", -1)
			cfg.Tracks[t] = strings.Replace(cfg.Tracks[t], "_", " ", -1)
		}
	})
	return cfg
}
