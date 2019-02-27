<h1 align="center">terraform-linter</h1>

<p align="center">
  Binary that lints a set of terraform files to adhere to the Vidsy linting guidelines.
</p>


[![Documentation](https://godoc.org/github.com/vidsy/terraform-linter?status.svg)](https://godoc.org/github.com/vidsy/terraform-linter)

![image](https://user-images.githubusercontent.com/527874/53488721-0a574c00-3a87-11e9-84f9-2245a505bc97.png)

# Vidsy linting

Terraform stacks can quickly get out of sync as different people and teams work on them. At Vidsy we have a set of simple guidelines for stacks to try and keep them consistent and easy to navigate and read.

# Usage

## Releases

The binary is versioned and released on each tagged merge to master, this can be found in the [releases](https://github.com/vidsy/terraform-linter/releases).

Once downloaded and installed, run the following to lint your stack:

```
terraform-linter --tf-directory"/path/to/terrform/files"
```

## Docker

The binary is also built to a container and pushed up to docker hub. To lint the files in the current directory run:

```
docker run --rm=true -v ${pwd}:/stack vidsyhq/terraform-linter --tf-directory="/stack"
```

## Linting rules

The following files are linted within the given stack:

### providers.tf

If this file exists, the following is checked:

1. Should only contain 1 or more `provider` configs and one `terraform` config.
1. Should contain no `data`, `local`, `module`, `output` or`resource` resources.

### resources.tf

If this file exists, the following is checked:

1. Should contain 1 or more `local, `module` or `resource` resources.
1. Should not contain `data`, `provider`, `terraform` or `output` resources.

### data.tf

If this file exists, the following is checked:

1. Should contain 1 or more `data` resources.
1. Should not contain `local`, `module`, `output`, `provider`, `resource` or `terraform` resources.

### outputs.tf

If this file exists, the following is checked:

1. Should contain 1 or more `output` resources.
1. Should not contain `data`, `local`, `module`, `provider`, `resource` or `terraform` resources.
