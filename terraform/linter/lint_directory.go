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
		config, err := config.LoadFile(filePath)
		if err != nil {
			return errors.Wrapf(
				err,
				"Problem parsing terraform config in %s",
				filePath,
			)
		}

		var linters []Linter

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
			if err := linter(config); err != nil {
				return NewError(
					err,
					file.Name(),
				)
			}
		}
	}

	return nil

}

func isValidTFFile(file os.FileInfo) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == ".tf"
}

func shouldNotContain(conf *config.Config, types ...terraformType) error {
	var err error
	errorMessage := "should not contain any %s resource(s), please move to '%s.tf'"

	for _, terraformType := range types {
		switch terraformType {
		case terraformTypeData:
			for _, resource := range conf.Resources {
				if resource.Mode == config.DataResourceMode {
					return errors.Errorf(errorMessage, "data", "data")
				}
			}
		case terraformTypeResource:
			for _, resource := range conf.Resources {
				if resource.Mode == config.ManagedResourceMode {
					return errors.Errorf(errorMessage, "resource", "resources")
				}
			}
		case terraformTypeOutput:
			if len(conf.Outputs) > 0 {
				return errors.Errorf(errorMessage, "output", "outputs")
			}
		case terraformTypeVariable:
			if len(conf.Variables) > 0 {
				return errors.Errorf(errorMessage, "variable", "variables")
			}
		case terraformTypeLocal:
			if len(conf.Locals) > 0 {
				return errors.Errorf(errorMessage, "local", "resources")
			}
		case terraformTypeProvider:
			if len(conf.ProviderConfigs) > 0 {
				return errors.Errorf(errorMessage, "provider", "providers")
			}

		case terraformTypeTerraform:
			if conf.Terraform != nil {
				return errors.Errorf(errorMessage, "terraform", "providers")
			}
		}
	}

	return err
}
