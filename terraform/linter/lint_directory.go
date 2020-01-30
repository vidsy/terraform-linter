package linter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

// LintDirectory takes an set files an lints them based on the Vidsy
// structure of stacks.
func LintDirectory(directory string, files []os.FileInfo) error {
	for _, info := range files {
		if !isValidTFFile(info) {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", directory, info.Name())

		parts := strings.Split(info.Name(), ".")
		name := parts[0]

		tf, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%s*.tf", name))
		if err != nil {
			return err
		}

		defer tf.Close()

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		defer file.Close()

		buf := bytes.NewBuffer([]byte{})
		if _, err := buf.Write([]byte("#terraform:hcl2\n")); err != nil {
			return err
		}
		if _, err := buf.ReadFrom(file); err != nil {
			return err
		}

		if _, err := tf.Write(buf.Bytes()); err != nil {
			return err
		}

		fmt.Println(tf.Name())

		defer os.Remove(tf.Name())

		conf, err := config.LoadFile(tf.Name())
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
			linters = append(linters, LintResources)
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
