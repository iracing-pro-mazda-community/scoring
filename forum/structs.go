package forum

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Post struct {
	ID      string
	Name    string
	Message string
}

func (p *Post) Store() error {
	filename := filepath.Join("output", fmt.Sprintf("%s - %s.txt", p.Name, p.ID))
	return ioutil.WriteFile(filename, []byte(p.Message), 0644)
}
