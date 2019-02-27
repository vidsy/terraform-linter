package linter

import (
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
