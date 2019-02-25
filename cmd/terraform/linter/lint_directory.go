package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/hashicorp/terraform/config"
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

		switch file.Name() {
		case "data.tf":
			err = lintData(config)
		case "providers.tf":
			err = lintProviders(config)
		case "resources.tf":
			err = lintResources(config)
		case "variables.tf":
			err = lintVariables(config)
		}

		if err != nil {
			return errors.Errorf(
				"%s - %s",
				file.Name(),
				err,
			)
		}
	}

	return nil

}

func lintData(terraformConfig *config.Config) error {
	if len(terraformConfig.Resources) == 0 {
		return errors.New(
			"contains no resources, either remove the file or add some data resources",
		)
	}

	for _, resource := range terraformConfig.Resources {
		if resource.Mode != config.DataResourceMode {
			return errors.New(
				"should only contain data resources, please remove",
			)
		}
	}

	return nil
}

func lintProviders(terraformConfig *config.Config) error {
	if terraformConfig.Terraform == nil {
		return errors.New(
			"contains no terraform resource, either remove the file or add a terraform resource",
		)
	}

	if len(terraformConfig.ProviderConfigs) == 0 {
		return errors.New(
			"contains no provider resources, if other resources exist add one or remove the file",
		)
	}

	if len(terraformConfig.Resources) > 0 {
		return errors.Errorf(
			"contains %d resources, please move to either 'resources.tf' or 'data.tf' depending on type",
			len(terraformConfig.Resources),
		)
	}

	return nil
}

func lintResources(terraformConfig *config.Config) error {
	if terraformConfig.Terraform != nil {
		return errors.New(
			"contains a terraform resource, this should be placed in 'providers.tf'",
		)
	}

	if len(terraformConfig.ProviderConfigs) > 0 {
		return errors.Errorf(
			"contains %d provider resource(s), these should be placed in 'providers.tf'",
			len(terraformConfig.ProviderConfigs),
		)
	}

	if len(terraformConfig.Resources) == 0 &&
		len(terraformConfig.Modules) == 0 &&
		len(terraformConfig.Locals) == 0 {
		return errors.New(
			"contains no resources, modules or locals. Either remove the file or add some resources/modules/locals",
		)
	}

	for _, resource := range terraformConfig.Resources {
		if resource.Mode == config.DataResourceMode {
			return errors.New(
				"should not contain any data resources, please move to 'data.tf'",
			)
		}
	}

	return nil
}

func lintVariables(terraformConfig *config.Config) error {
	if len(terraformConfig.Variables) == 0 {
		return errors.New(
			"no variables found, either add some or remove the file",
		)
	}

	for _, variable := range terraformConfig.Variables {
		if val, ok := variable.Default.(string); ok && val == "" {
			return errors.Errorf(
				"variable '%s' contains a blank default, please remove the default",
				variable.Name,
			)
		}
	}

	return nil
}

func isValidTFFile(file os.FileInfo) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == ".tf"
}
