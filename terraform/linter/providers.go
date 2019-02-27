package linter

import (
	"github.com/hashicorp/terraform/config"
	"github.com/pkg/errors"
)

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
