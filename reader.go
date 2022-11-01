package main

import (
	"io"
	"io/ioutil"
	"os"
)

func ReadFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

func ReadFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
