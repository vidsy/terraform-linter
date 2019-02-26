package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/vidsy/terraform-linter/terraform/linter"
)

func main() {
	tfDirectory := flag.String(
		"tf-directory",
		"",
		"The directory that contains the terraform files to lint",
	)

	hideStackTraces := flag.Bool(
		"hide-stack-traces",
		true,
		"Should stack traces be shown for errors",
	)

	flag.Parse()

	err := isValidDirectory(*tfDirectory)
	if err != nil {
		linter.PrintFatalError(err, *hideStackTraces)
	}

	files, err := ioutil.ReadDir(*tfDirectory)
	if err != nil {
		linter.PrintFatalError(err, *hideStackTraces)
	}

	err = linter.LintDirectory(*tfDirectory, files)
	if err != nil {
		linter.PrintFatalError(err, *hideStackTraces)
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
