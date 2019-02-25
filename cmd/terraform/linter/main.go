package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

func main() {
	tfDirectory := flag.String(
		"tf-directory",
		"",
		"The directory that contains the terraform files to lint",
	)
	flag.Parse()

	err := isValidDirectory(*tfDirectory)
	if err != nil {
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

func isValidDirectory(path string) error {
	if path == "" {
		return errors.New("tf-directory is blank")
	}

	directoryInfo, err := os.Stat(path)
	if err != nil {
		return errors.Wrap(err, "Problem stating directory")
	}

	if !directoryInfo.IsDir() {
		return errors.Errorf(
			"Expected '%s' to be a directory",
			path,
		)
	}

	return nil
}
