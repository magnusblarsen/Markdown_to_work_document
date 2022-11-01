package main 

import(
    "strings"
    "os"
    "bufio"
    "log"
)

func ParseToMarkdown(fileName string) ([]byte) {
	file, err := os.Open(fileName)
    checkError(err)

    fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)

    ScanVariables(fileScanner)
    parsedData := ReplaceVariables(fileScanner)

    file.Close()

    return parsedData
}

func ScanVariables(fileScanner *bufio.Scanner) {
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

func ReplaceVariables(fileScanner *bufio.Scanner) []byte {
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
