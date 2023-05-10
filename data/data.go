package data

import (
	"os"
)

type Data struct {
	content map[string]string
}

func NewData(files ...string) (*Data, error) {
	data := &Data{content: map[string]string{}}
	for _, name := range files {
		b, err := os.ReadFile("data/content/" + name + ".txt")
		if err != nil {
			return nil, err
		}
		data.content[name] = string(b)
	}

	return data, nil
}

func (d *Data) Content(name string) string {
	return d.content[name]
}
