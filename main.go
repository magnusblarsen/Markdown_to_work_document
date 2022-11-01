package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gitlab.com/golang-commonmark/markdown"
)

var buf []byte
var variables = make(map[string]string)

func main() {
	if len(os.Args) < 2 {
		log.Println("You need to specify a working path")
		os.Exit(1)
	}
	buf = make([]byte, 128) //TODO: find right buffer length

	dir := os.Args[1]
	regex, err := regexp.Compile(".*.md$")
	checkError(err)
	filePathErr := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		if !regex.MatchString(path) {
			return err
		}
		outputFile := ParseToMarkdown(path)
        fmt.Println(convertToHtml(outputFile)) //TODO: Write to file
		return err
	})
	checkError(filePathErr)
}

func convertToHtml(markdownFile []byte) string {
	md := markdown.New(markdown.HTML(true))
	return md.RenderToString(markdownFile)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
