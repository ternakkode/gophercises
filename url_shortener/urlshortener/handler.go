package urlshortener

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"title"`
	URL  string `yaml:"url"`
}

func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathToUrls[path]; ok {
			http.Redirect(rw, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(rw, r)
	}
}

func YamlHandler(yamlByte []byte, fallback http.Handler) http.HandlerFunc {
	pathUrl, err := parseYaml(yamlByte)
	if err != nil {
		log.Fatalln(err)
	}
	pathUrlsMap := buildMap(pathUrl)
	return MapHandler(pathUrlsMap, fallback)
}

func parseYaml(yamlByte []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl

	err := yaml.Unmarshal(yamlByte, &pathUrls)
	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathUrlsMap := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathUrlsMap[pathUrl.Path] = pathUrl.URL
	}

	return pathUrlsMap
}
