package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"quiet_hn/hn"
)

var (
	cache           []item
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getCachedStories(numStories)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}

		err = tpl.Execute(w, data)

		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getCachedStories(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().Before(cacheExpiration) {
		return cache, nil
	}

	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}

	cache = stories
	cacheExpiration = time.Now().Add(1 * time.Minute)

	return stories, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()

	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	resultCh := make(chan ItemResult)
	for i := 0; i <= numStories*5/4; i++ {
		go getStoryAsync(ids[i], resultCh)
	}

	var items []item
	for i := 0; i < numStories*5/4; i++ {
		result := <-resultCh

		if isStoryLink(result.item) {
			items = append(items, result.item)
			if len(items) == numStories {
				break
			}
		}
	}

	return items, nil
}

func getStoryAsync(id int, channel chan ItemResult) {
	var client hn.Client
	hnItem, err := client.GetItem(id)
	if err != nil {
		channel <- ItemResult{err: err}
	}

	item := parseHNItem(hnItem)
	channel <- ItemResult{item: item, index: id}
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

type ItemResult struct {
	item  item
	index int
	err   error
}
