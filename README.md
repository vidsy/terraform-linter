<h1 align="center">terraform-liner</h1>

<p align="center">
  Binary that lints a set of terraform files to adhere to the Vidsy linting guidelines.
</p>


[![Documentation](https://godoc.org/github.com/vidsy/terraform-linter?status.svg)](https://godoc.org/github.com/vidsy/terraform-linter)

# Vidsy linting

Terraform stacks can quickly get out of sync as different people and teams work on them. At Vidsy we have a set of simple guidelines for stacks to try and keep them consistent and easy to navigate and read.

## Linting rules

The following files are checked for certain linting requirements:

### providers.tf

If this file exists, the following is checked:

1. Should only contain 1 or more `provider` configs and one `terraform` config.
1. Should contain no `modules`, `resources`, `outputs` or `data` definitions.

## resources.tf

If this file exists, the following is checked:

1. Should contain 1 or more `resource`, `module` or `local` definitions.
1. Should not contain `provider`, `terraform` or `output` definitions.

## data.tf

If this file exists, the following is checked:

1. Should contain 1 or more `data` definitions.
1. Should not contain `provider`, `resource`, `local`, `output` or `terraform` definitions.

## outputs.tf

If this file exists, the following is checked:

1. Should contain 1 or more `output` definitions.
1. Should not contain `provider`, `resource`, `local`, `data` or `terraform` definitions.
