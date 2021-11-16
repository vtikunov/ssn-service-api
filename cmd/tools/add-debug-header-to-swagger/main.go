package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please set a file name.")
		return
	}

	filename := os.Args[1]

	//nolint:gosec
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Can not open file: %s", err.Error())
		return
	}
	defer func() {
		if errCl := f.Close(); errCl != nil {
			log.Error().Err(errCl).Msg("failed closing file")
		}
	}()

	var lines []string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if lineTr := strings.Trim(line, " "); lineTr == `"parameters": [` {
			lines = append(lines, `          {
            "name": "Grpc-Metadata-Log-Level",
            "in": "header",
            "required": false,
            "type": "string",
            "format": "string"
          },`)
		}
	}

	if err = scanner.Err(); err != nil {
		log.Error().Err(err).Msg("read file error")
		return
	}

	//nolint:gosec
	err = ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)

	if err != nil {
		log.Error().Err(err).Msg("write file error")
	}

	return
}
