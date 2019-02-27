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

		switch file.Name() {
		case "data.tf":
			err = lintData(config)
		case "outputs.tf":
			err = lintOutputs(config)
		case "providers.tf":
			err = lintProviders(config)
		case "resources.tf":
			err = lintResources(config)
		case "variables.tf":
			err = lintVariables(config)
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

func isValidTFFile(file os.FileInfo) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == ".tf"
}

func lintData(conf *config.Config) error {
	err := shouldNotContain(
		conf,
		terraformTypeLocal,
		terraformTypeModule,
		terraformTypeOutput,
		terraformTypeProvider,
		terraformTypeResource,
		terraformTypeTerraform,
		terraformTypeVariable,
	)
	if err != nil {
		return err
	}

	if len(conf.Resources) == 0 {
		return errors.New(
			"contains no data resources, either remove the file or add some data resources",
		)
	}

	for _, resource := range conf.Resources {
		if resource.Mode != config.DataResourceMode {
			return errors.New(
				"should only contain data resources, please remove",
			)
		}
	}

	return nil
}

func lintOutputs(conf *config.Config) error {
	err := shouldNotContain(
		conf,
		terraformTypeData,
		terraformTypeLocal,
		terraformTypeModule,
		terraformTypeProvider,
		terraformTypeResource,
		terraformTypeTerraform,
		terraformTypeVariable,
	)
	if err != nil {
		return err
	}

	if len(conf.Outputs) == 0 {
		return errors.New(
			"no outputs found, either add some or remove the file",
		)
	}

	return nil
}

func lintProviders(conf *config.Config) error {
	err := shouldNotContain(
		conf,
		terraformTypeData,
		terraformTypeLocal,
		terraformTypeModule,
		terraformTypeOutput,
		terraformTypeResource,
		terraformTypeVariable,
	)
	if err != nil {
		return err
	}

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

	return nil
}

func lintResources(conf *config.Config) error {
	err := shouldNotContain(
		conf,
		terraformTypeData,
		terraformTypeOutput,
		terraformTypeProvider,
		terraformTypeTerraform,
		terraformTypeVariable,
	)
	if err != nil {
		return err
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

func lintVariables(conf *config.Config) error {
	err := shouldNotContain(
		conf,
		terraformTypeData,
		terraformTypeLocal,
		terraformTypeModule,
		terraformTypeOutput,
		terraformTypeProvider,
		terraformTypeResource,
		terraformTypeTerraform,
	)
	if err != nil {
		return err
	}

	if len(conf.Variables) == 0 {
		return errors.New(
			"no variables found, either add some or remove the file",
		)
	}

	for _, variable := range conf.Variables {
		if val, ok := variable.Default.(string); ok && val == "" {
			return errors.Errorf(
				"variable '%s' contains a blank default, please remove the default",
				variable.Name,
			)
		}
	}

	return nil
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
