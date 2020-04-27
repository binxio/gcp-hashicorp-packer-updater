package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func updatePackerDefinition(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return err
	}

	var dirty bool
	switch builders := result["builders"].(type) {
	case []interface{}:
		for i, b := range builders {
			switch builder := b.(type) {
			case map[string]interface{}:
				builder_content, err := json.Marshal(b)
				if err != nil {
					return err
				}
				var gceBuilder GoogleComputeBuilder
				err = json.Unmarshal(builder_content, &gceBuilder)
				if err != nil {
					return err

				}
				if gceBuilder.Type == "googlecompute" {
					updated, err := gceBuilder.updateGoogleSourceImage()
					if err != nil {
						return err
					}
					if updated {
						builder["source_image"] = gceBuilder.Image
						builder["source_image_family"] = gceBuilder.Family
						builder["source_image_project_id"] = gceBuilder.ProjectId
					}
					dirty = dirty || updated
				}
			default:
				return fmt.Errorf("builder[%d] in %s is not an object", i, filename)
			}
		}
	default:
		return fmt.Errorf("builders in %s is not an array of objects", filename)
	}

	if dirty {
		x, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Errorf("failed to marshal updated file to json, %s\n", err)
		}
		ioutil.WriteFile(filename, x, os.FileMode(0660))
	}
	return nil
}

func main() {
	var filename = flag.String("filename", "packer.json", "of the packer file to update")

	flag.Parse()
	err := updatePackerDefinition(*filename)
	if err != nil {
		log.Fatal(err)
	}
}
