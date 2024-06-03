package sysdlux

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alexander-bergeron/go-systemdelux/monitor"
	"gopkg.in/yaml.v2"
)

func AddCommand(args []string) {
	var yamlFile string
	var yamlString string

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	// TODO: Expand to json as well
	fs.StringVar(&yamlFile, "file", "", "Path to the YAML file")
	fs.StringVar(&yamlString, "yaml", "", "YAML string")

	fs.Parse(args)

	var svc monitor.Service

	// TODO: Add better validation
	if yamlFile != "" {
		data, err := os.ReadFile(yamlFile)
		if err != nil {
			log.Fatalf("Failed to read YAML file: %v", err)
		}
		err = yaml.Unmarshal(data, &svc)
		if err != nil {
			log.Fatalf("Failed to unmarshal YAML file: %v", err)
		}
	} else if yamlString != "" {
		err := yaml.Unmarshal([]byte(yamlString), &svc)
		if err != nil {
			log.Fatalf("Failed to unmarshal YAML string: %v", err)
		}
	} else {
		log.Fatalf("Either a YAML file path or a YAML string must be provided")
	}

	// Marshal the struct into a YAML byte slice
	data, err := yaml.Marshal(svc)
	if err != nil {
		log.Fatalf("failed to marshal struct: %v", err)
	}
	// TODO: do with path module
	// TODO: check for conflicts
	filename := fmt.Sprintf("%v/%v.yml", svcPath, svc.Name)

	// Write the YAML byte slice to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}
}
