package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Configuration struct {
	Topic    string
	Tracks   []string
	Download bool
}

func Get() (*Configuration, error) {
	file, err := os.Open("configuration.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	if err := decoder.Decode(&configuration); err != nil {
		return nil, err
	}

	// lowercase track names
	for t := range configuration.Tracks {
		configuration.Tracks[t] = strings.ToLower(configuration.Tracks[t])
	}

	return &configuration, nil
}
