package linter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			err = lintData(config.Resources)
		case "providers.tf":
			err = lintProviders(config)
		case "resources.tf":
			err = lintResources(config)
		case "variables.tf":
			err = lintVariables(config.Variables)
		}

		if err != nil {
			return NewError(
				err,
				file.Name(),
			)
		}
	}

	return nil
}

func lintName(conf *config.Config) error {
	for _, variable := range conf.Variables {
		if !isValidName(variable.Name) {
			return errors.Errorf(
				"variable name '%s' contains hyphens, please replace with underscores",
				variable.Name,
			)
		}
	}

	for _, resource := range conf.Resources {
		if !isValidName(resource.Name) {
			return errors.Errorf(
				"resource name '%s' contains hyphens, please replace with underscores",
				resource.Name,
			)
		}
	}

	for _, module := range conf.Modules {
		if !isValidName(module.Name) {
			return errors.Errorf(
				"module name '%s' contains hyphens, please replace with underscores",
				module.Name,
			)
		}
	}

	for _, local := range conf.Locals {
		if !isvalidname(local.Name) {
			return errors.errorf(
				"local name '%s' contains hyphens, please replace with underscores",
				local.Name,
			)
		}
	}

	for _, output := range conf.Outputs {
		if !isvalidname(output.Name) {
			return errors.errorf(
				"output name '%s' contains hyphens, please replace with underscores",
				output.Name,
			)
		}
	}
}

func isValidName(name string) bool {
	return strings.Contains(name, "-")
}

func lintData(resources []*config.Resource) error {
	if len(resources) == 0 {
		return errors.New(
			"contains no resources, either remove the file or add some data resources",
		)
	}

	for _, resource := range resources {
		if resource.Mode != config.DataResourceMode {
			return errors.New(
				"should only contain data resources, please remove",
			)
		}
	}

	return nil
}

func lintProviders(conf *config.Config) error {
	if conf.Terraform == nil {
		return errors.New(
			"contains no terraform resource, either remove the file or add a terraform resource",
		)
	}

	if len(conf.ProviderConfigs) == 0 {
		return errors.New(
			"contains no provider resources, if other resources exist add one or remove the file",
		)
	}

	if len(conf.Resources) > 0 {
		return errors.Errorf(
			"contains %d resources, please move to either 'resources.tf' or 'data.tf' depending on type",
			len(conf.Resources),
		)
	}

	return nil
}

func lintResources(conf *config.Config) error {
	if conf.Terraform != nil {
		return errors.New(
			"contains a terraform resource, this should be placed in 'providers.tf'",
		)
	}

	if len(conf.ProviderConfigs) > 0 {
		return errors.Errorf(
			"contains %d provider resource(s), these should be placed in 'providers.tf'",
			len(conf.ProviderConfigs),
		)
	}

	if len(conf.Resources) == 0 &&
		len(conf.Modules) == 0 &&
		len(conf.Locals) == 0 {
		return errors.New(
			"contains no resources, modules or locals. Either remove the file or add some resources/modules/locals",
		)
	}

	for _, resource := range conf.Resources {
		if resource.Mode == config.DataResourceMode {
			return errors.New(
				"should not contain any data resources, please move to 'data.tf'",
			)
		}
	}

	return nil
}

func lintVariables(variables []*config.Variable) error {
	if len(variables) == 0 {
		return errors.New(
			"no variables found, either add some or remove the file",
		)
	}

	for _, variable := range variables {
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
