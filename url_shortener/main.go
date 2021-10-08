package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"url_shortener/urlshortener"
)

// problem statement :
// make http server with redirection to specific url
func main() {
	mux := defaultMux()

	var handler interface{}
	if file, err := os.ReadFile("url.yaml"); err != nil {
		log.Println("failed to parse yaml file, existing data used...")
		handler = urlshortener.MapHandler(getMapPathToUrl(), mux)
	} else {
		handler = urlshortener.YamlHandler(file, mux)
	}

	fmt.Println("http server started on port 8080")
	http.ListenAndServe(":8080", handler.(http.HandlerFunc))
}

func getMapPathToUrl() map[string]string {
	return map[string]string{
		"/rabraw": "https://rajabrawijaya.ub.ac.id",
		"/blibli": "https://blibli.com",
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
