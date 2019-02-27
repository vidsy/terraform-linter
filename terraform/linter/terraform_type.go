package linter

const (
	terraformTypeData      terraformType = "data"
	terraformTypeLocal     terraformType = "local"
	terraformTypeModule    terraformType = "module"
	terraformTypeOutput    terraformType = "output"
	terraformTypeProvider  terraformType = "provider"
	terraformTypeResource  terraformType = "resource"
	terraformTypeTerraform terraformType = "terraform"
	terraformTypeVariable  terraformType = "variable"
)

type (
	terraformType string
)
