package web

import (
	"adventure_book/model"
	"adventure_book/reader"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	htmlTemplate, err := reader.GetHtmlTemplate()
	if err != nil {
		log.Panicln(err)
	}

	tpl = template.Must(template.New("").Parse(htmlTemplate))
}

func Start(stories model.Story) {
	handler := NewHandler(stories)

	fmt.Println("starting server on port :8080")
	http.ListenAndServe(":8080", handler)
}

type Handler struct {
	s model.Story
}

func NewHandler(s model.Story) http.Handler {
	return Handler{s: s}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	chapter := path[1:]
	if chapterDetail, ok := h.s[chapter]; ok {
		err := tpl.Execute(w, chapterDetail)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)
}
