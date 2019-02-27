package linter

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

const (
	validNameRegex = `^[a-z0-9_]+$`
)

// LintData lints resources in a data.tf file.
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
// data and output resources for dashes.
func LintNames(conf *config.Config) error {
	errorMessage := fmt.Sprintf(
		"%%s name '%%s' is not valid for the regex '%s', please replace the invalid characters",
		validNameRegex,
	)

	for _, local := range conf.Locals {
		if !isValidName(local.Name) {
			return errors.Errorf(errorMessage, "local", local.Name)
		}
	}

	for _, module := range conf.Modules {
		if !isValidName(module.Name) {
			return errors.Errorf(errorMessage, "module", module.Name)
		}
	}

	for _, output := range conf.Outputs {
		if !isValidName(output.Name) {
			return errors.Errorf(errorMessage, "output", output.Name)
		}
	}

	for _, provider := range conf.ProviderConfigs {
		if !isValidName(provider.Name) {
			return errors.Errorf(errorMessage, "provider", provider.Name)
		}
	}

	for _, resource := range conf.Resources {
		if !isValidName(resource.Name) {
			switch resource.Mode {
			case config.DataResourceMode:
				return errors.Errorf(errorMessage, "data", resource.Name)
			case config.ManagedResourceMode:
				return errors.Errorf(errorMessage, "resource", resource.Name)
			}
		}
	}

	for _, variable := range conf.Variables {
		if !isValidName(variable.Name) {
			return errors.Errorf(errorMessage, "variable", variable.Name)
		}
	}

	return nil
}

// LintOutputs lints outputs in a outputs.tf file.
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
// types in a providers.tf file.
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

// LintResources lints resources, modules and locals in
// a resources.tf file.
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

	return nil
}

// LintVariables lints variables in a variables.tf file.
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
	nameRegex := regexp.MustCompile(validNameRegex)
	return nameRegex.MatchString(name)
}

func shouldNotContain(conf *config.Config, types ...terraformType) error {
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

	return nil
}
