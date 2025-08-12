package balancer

import (
	"os"
	"strings"
)

func ReadURLs() ([]string, error) {
	file, err := os.Open("./data/urls.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 4096)
	n, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	urls := strings.Split(string(data[:n]), ",")
	return urls, nil
}

func WriteURL(url string) error {
    file, err := os.OpenFile("./data/urls.csv", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(url + ",")
    return err
}