package linter

import (
	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

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
