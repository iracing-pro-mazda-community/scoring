package forum

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/iracing-pro-mazda-community/scoring/config"
)

type Post struct {
	ID      string
	Name    string
	Message string
}

func (p *Post) Store() error {
	filename := filepath.Join("output", fmt.Sprintf("%s - %s.txt", p.Name, p.ID))

	if !config.Get().OverwriteVote {
		if _, err := os.Stat(filename); err == nil {
			log.Printf("[%s] already exists, not overwriting!\n", filename)
			return nil
		}
	}

	return ioutil.WriteFile(filename, []byte(p.Message), 0644)
}
