package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/config"
)

func main() {
	tfDirectory := flag.String("tf-directory", "", "The directory that contains the terraform files to lint")
	flag.Parse()

	if *tfDirectory == "" {
		log.Fatal(
			errors.New("tf-directory is blank"),
		)
	}

	directoryInfo, err := os.Stat(*tfDirectory)
	if err != nil {
		log.Fatal(err)
	}

	if !directoryInfo.IsDir() {
		log.Fatal(
			fmt.Errorf(
				"Expected '%s' to be a directory",
				*tfDirectory,
			),
		)
	}

	files, err := ioutil.ReadDir(*tfDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".tf" {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", *tfDirectory, file.Name())
		config, err := config.LoadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		switch file.Name() {
		case "resources.tf":
			err = lintResources(config)
		case "data.tf":
			err = lintData(config)
		case "providers.tf":
			err = lintProviders(config)
		}

		if err != nil {
			log.Fatalf("%s - %s", file.Name(), err)
		}

		log.Printf("Linting complete for '%s', no issues", file.Name())
	}

	log.Println("Linting finished")
}

func lintData(terraformConfig *config.Config) error {
	if len(terraformConfig.Resources) == 0 {
		return errors.New(
			"contains no resources, either remove the file or add some data resources",
		)
	}

	for _, resource := range terraformConfig.Resources {
		if resource.Mode != config.DataResourceMode {
			return errors.New(
				"should only contain data resources, please remove",
			)
		}
	}

	return nil
}

func lintProviders(terraformConfig *config.Config) error {
	if terraformConfig.Terraform == nil {
		return errors.New(
			"contains no terraform resource, either remove the file or add a terraform resource",
		)
	}

	if len(terraformConfig.ProviderConfigs) == 0 {
		return errors.New(
			"contains no provider resources, if other resources exist add one or remove the file",
		)
	}

	if len(terraformConfig.Resources) > 0 {
		return fmt.Errorf(
			"contains %d resources, please move to either 'resources.tf' or 'data.tf' depending on type",
			len(terraformConfig.Resources),
		)
	}

	return nil
}

func lintResources(terraformConfig *config.Config) error {
	if terraformConfig.Terraform != nil {
		return errors.New(
			"contains a terraform resource, this should be placed in 'providers.tf'",
		)
	}

	if len(terraformConfig.ProviderConfigs) > 0 {
		return fmt.Errorf(
			"contains %d provider resource(s), these should be placed in 'providers.tf'",
			len(terraformConfig.ProviderConfigs),
		)
	}

	if len(terraformConfig.Resources) == 0 {
		return errors.New(
			"contains no resources, either remove the file or add some resources",
		)
	}

	for _, resource := range terraformConfig.Resources {
		if resource.Mode == config.DataResourceMode {
			return errors.New(
				"should not contain any data resources, please move to 'data.tf'",
			)
		}
	}

	return nil
}
