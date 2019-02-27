package linter

import (
	"strings"

	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

// LintData lints data.
func LintData(conf *config.Config) error {
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

// LintNames linters the names of resources, modules, providers, locals, variables,
// data and output resources.
func LintNames(conf *config.Config) error {
	for _, local := range conf.Locals {
		if !isValidName(local.Name) {
			return errors.Errorf(
				"local name '%s' contains hyphens, please replace with underscores",
				local.Name,
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

	for _, output := range conf.Outputs {
		if !isValidName(output.Name) {
			return errors.Errorf(
				"output name '%s' contains hyphens, please replace with underscores",
				output.Name,
			)
		}
	}

	for _, provider := range conf.ProviderConfigs {
		if !isValidName(provider.Name) {
			return errors.Errorf(
				"provider name '%s' contains hyphens, please replace with underscores",
				provider.Name,
			)
		}
	}

	for _, resource := range conf.Resources {
		if !isValidName(resource.Name) {
			errorMessage := "%s name '%s' contains hyphens, please replace with underscores"

			switch resource.Mode {
			case config.DataResourceMode:
				return errors.Errorf(
					errorMessage,
					"data",
					resource.Name,
				)

			case config.ManagedResourceMode:
				return errors.Errorf(
					errorMessage,
					"resource",
					resource.Name,
				)
			}
		}
	}

	for _, variable := range conf.Variables {
		if !isValidName(variable.Name) {
			return errors.Errorf(
				"variable name '%s' contains hyphens, please replace with underscores",
				variable.Name,
			)
		}
	}

	return nil
}

// LintOutputs lints outputs.
func LintOutputs(conf *config.Config) error {
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

// LintProviders lints providers and terraform
// types.
func LintProviders(conf *config.Config) error {
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

// LintResources lints resources.
func LintResources(conf *config.Config) error {
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

// LintVariables lints variables.
func LintVariables(conf *config.Config) error {
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

func isValidName(name string) bool {
	return !strings.Contains(name, "-")
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
