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
	tfDirectory := flag.String("tf-directory", "", "The directory that contains the terraform files to lint")
	flag.Parse()

	if *tfDirectory == "" {
		log.Fatal(
			errors.New("tf-directory is blank"),
		)
	}

	directoryInfo, err := os.Stat(*tfDirectory)
	if err != nil {
		log.Fatal(err)
	}

	if !directoryInfo.IsDir() {
		log.Fatal(
			fmt.Errorf(
				"Expected '%s' to be a directory",
				*tfDirectory,
			),
		)
	}

	files, err := ioutil.ReadDir(*tfDirectory)
	if err != nil {
		log.Fatal(err)
	}

	err = Linter(*tfDirectory, files)
	if err != nil {
		log.Fatal(err)
	}
}
