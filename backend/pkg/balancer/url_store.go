package balancer

import (
	"os"
	"strings"
)

const URL_STORE_FILEPATH string = "./data/urls.csv"

/*
These functions are intended to be used to create / update
a permanent URL store
*/


func ReadURLs() ([]string, error) {
	file, err := os.Open(URL_STORE_FILEPATH)
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
	file, err := os.OpenFile(URL_STORE_FILEPATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url + ",")
	return err
}
