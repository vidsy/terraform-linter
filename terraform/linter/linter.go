package linter

import "github.com/hashicorp/terraform/config"

type (
	// Linter type for functions that can lint a terraform config
	// object and return an error if there's an issue.
	Linter func(*config.Config) error
)
