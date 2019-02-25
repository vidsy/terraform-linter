package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	tfDirectory := flag.String(
		"tf-directory",
		"",
		"The directory that contains the terraform files to lint",
	)
	flag.Parse()

	if valid, err := isValidDirectory(*tfDirectory); !valid {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(*tfDirectory)
	if err != nil {
		log.Fatal(err)
	}

	err = LintDirectory(*tfDirectory, files)
	if err != nil {
		log.Fatal(err)
	}
}

func isValidDirectory(path string) (bool, error) {
	if path == "" {
		return false, errors.New("tf-directory is blank")
	}

	directoryInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !directoryInfo.IsDir() {
		return false, fmt.Errorf(
			"Expected '%s' to be a directory",
			path,
		)
	}

	return true, nil
}
