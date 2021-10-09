package reader

import "os"

func GetHtmlTemplate() (string, error) {
	file, err := os.ReadFile("./static/template/default.html")
	if err != nil {
		return "", err
	}

	return string(file), nil
}
