package main

import (
	"adventure_book/reader"
	"adventure_book/web"
	"log"
)

// problem statement :
// given a story
// generate html to show the story

func main() {
	stories, err := reader.ReadJsonStory("./static/story/default.json")
	if err != nil {
		log.Panicln(err)
	}

	web.Start(stories)
}
