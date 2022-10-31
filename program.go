package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

var buf []byte
var variables = make(map[string]string)


func main () {
    if len(os.Args) < 2 {
        log.Println("You need to specify a working path")
        os.Exit(1)
    }
    buf = make([]byte, 128) //TODO: find right buffer length

    dir := os.Args[1]
    fileList := make([]string, 0)
    regex, err := regexp.Compile(".*.md$")
    checkError(err)
    filePathErr := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if info.IsDir() {
            return err
        }
        if !regex.MatchString(path) {
            return err
        }
        fileList = append(fileList, path)

        outputFile := parseToMarkdown(path)
        fmt.Println(convertToHtml(outputFile))
        return err
    })
    checkError(filePathErr)
}

func parseToMarkdown(fileName string) ([]byte) {
	file, err := os.Open(fileName)
    checkError(err)

    fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)

    scanVariables(fileScanner)
    parsedData := replaceVariables(fileScanner)

    file.Close()

    return parsedData
}

func scanVariables(fileScanner *bufio.Scanner) {
    for fileScanner.Scan() {
        line := fileScanner.Text()

        var key []rune
        var value []rune

        parsingKey := true
        for _, v := range line {
            if v == '\\' {
                log.Println("Done parsing variables")
                return
            }
            if parsingKey == true && v == '=' {
                parsingKey = false
                continue
            }

            if parsingKey {
                key = append(key, v)
            } else{
                value = append(value, v)
            }
        }

        if parsingKey {
            log.Printf("Skipping line: %s\nNo variables found", line)
            continue
        }


        cutset := " "
        trimmedKey := strings.Trim(string(key), cutset)
        trimmedValue := strings.Trim(string(value), cutset)
        variables[trimmedKey] = trimmedValue
        log.Printf("key: %s with value: %s", trimmedKey, variables[trimmedKey])
    }
}

func replaceVariables(fileScanner *bufio.Scanner) []byte {
    parsedData := make([]byte, 0)
    for fileScanner.Scan() {
        line := fileScanner.Text()

        parsingVariable := false
        parsedVariable := make([]rune, 0)

        replaceVariable := func() {
            parsingVariable = false
            completeVariable := string(parsedVariable)
            log.Printf("replacing variable: %s with value: %s", completeVariable, variables[completeVariable])
            parsedData = append(parsedData, []byte(variables[string(parsedVariable)])...)
            parsedVariable = make([]rune, 0)
        }
        for i, v := range line {
            if v == '$' {
                parsingVariable = true
                continue
            }
            if parsingVariable == true {
                if v == ';' {
                    replaceVariable()
                } else if v == ' ' || v == ',' {
                    replaceVariable()
                    parsedData = append(parsedData, []byte(string(v))...)
                } else if i == len(line) - 1 {
                    parsedVariable = append(parsedVariable, v)
                    replaceVariable()
                }else{
                    parsedVariable = append(parsedVariable, v)
                }
                continue
            }
            parsedData = append(parsedData, []byte(string(v))...)
        }
        parsedData = append(parsedData, []byte("\n")...)
    }

    return parsedData
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

func readFromStdin () ([]byte, error) {
    return ioutil.ReadAll(os.Stdin)
}

func readFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

