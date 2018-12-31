package util

import (
	"io/ioutil"
	"os"
)

func ReadFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

func WriteFile(data string, file string) (int, error) {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return -1, err
	}

	n, err := f.WriteString(data)
	return n, err
}
