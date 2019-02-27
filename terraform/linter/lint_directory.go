package linter

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
			err = lintData(config.Resources)
		case "outputs.tf":
			err = lintOutputs(config.Outputs)
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

func lintOutputs(outputs []*config.Output) error {
	if len(outputs) == 0 {
		return errors.New(
			"no outputs found, either add some or remove the file",
		)
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
