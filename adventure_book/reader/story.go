package reader

import (
	"adventure_book/model"
	"encoding/json"
	"os"
)

func ReadJsonStory(filename string) (storyRes model.Story, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(file)
	if err := d.Decode(&storyRes); err != nil {
		return nil, err
	}

	return storyRes, nil
}
