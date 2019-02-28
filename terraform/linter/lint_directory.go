package linter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

// LintDirectory takes an set files an lints them based on the Vidsy
// structure of stacks.
func LintDirectory(directory string, files []os.FileInfo) error {
	for _, file := range files {
		if !isValidTFFile(file) {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", directory, file.Name())
		conf, err := config.LoadFile(filePath)
		if err != nil {
			return errors.Wrapf(
				err,
				"Problem parsing terraform config in %s",
				filePath,
			)
		}

		linters := []Linter{LintNames}

		switch file.Name() {
		case "data.tf":
			linters = append(linters, LintData)
		case "outputs.tf":
			linters = append(linters, LintOutputs)
		case "providers.tf":
			linters = append(linters, LintProviders)
		case "resources.tf":
			linters = append(linters, LintResources, LintUnusedVariables)
		case "variables.tf":
			linters = append(linters, LintVariables)
		}

		for _, linter := range linters {
			if err := linter(conf); err != nil {
				return NewError(
					err,
					file.Name(),
				)
			}
		}
	}

	conf, err := config.LoadDir(directory)
	if err != nil {
		return errors.Wrapf(
			err,
			"Problem parsing terraform config directory %s",
			directory,
		)
	}

	err = LintUnusedVariables(conf)
	if err != nil {
		return NewError(
			err,
			directory,
		)
	}

	return nil

}

func isValidTFFile(file os.FileInfo) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == ".tf"
}
