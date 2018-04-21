package score

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/iracing-pro-mazda-community/scoring/config"
)

var writer *csv.Writer

func init() {
	file, err := os.Create("output/result.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer = csv.NewWriter(file)
}

func WriteToCSV(data []string) {
	if err := writer.Write(data); err != nil {
		log.Fatal(err)
	}
	writer.Flush()
}

func WriteScoreToCSV(driver string, values map[string]int64) {
	data := []string{driver}
	for _, track := range config.Get().Tracks {
		if value, ok := values[track]; ok {
			data = append(data, fmt.Sprintf("%d", value))
		} else {
			data = append(data, fmt.Sprintf("%d", len(config.Get().Tracks)))
		}
	}
	WriteToCSV(data)
}
