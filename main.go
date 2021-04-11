package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/binxio/gcloudconfig"
	"golang.org/x/oauth2/google"
	"os"

	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type PackerUpdater struct {
	ctx         context.Context
	credentials *google.Credentials
	filename    string
}

func (updater *PackerUpdater) SourceImageDefinition() error {
	content, err := ioutil.ReadFile(updater.filename)
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
					updated, err := gceBuilder.updateGoogleSourceImage(updater.ctx, updater.credentials)
					if err != nil {
						return err
					}
					if updated {
						builder["source_image"] = gceBuilder.Image
						if gceBuilder.Family != "" {
							builder["source_image_family"] = gceBuilder.Family
						}
						if gceBuilder.ProjectId != "" {
							builder["source_image_project_id"] = gceBuilder.ProjectId
						}
					} else {
						log.Printf("builder[%d] image %s is up to date", i, gceBuilder.Image)
					}
					dirty = dirty || updated
				}
			default:
				return fmt.Errorf("builder[%d] in %s is not an object", i, updater.filename)
			}
		}
	default:
		return fmt.Errorf("builders in %s is not an array of objects", updater.filename)
	}

	if dirty {
		x, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Errorf("failed to marshal updated file to json, %s\n", err)
		}
		ioutil.WriteFile(updater.filename, x, os.FileMode(0660))
	}
	return nil
}

func main() {
	updater := PackerUpdater{ctx: context.Background()}

	var err error
	var project string
	useDefaultCredentials := flag.Bool("use-default-credentials", false, "for Google authentication")
	flag.StringVar(&updater.filename, "filename", "packer.json", "of the packer file to update")
	flag.StringVar(&project, "project", "", "to use ")
	configuration := flag.String("configuration", "", "of gcloud to use")
	flag.Parse()

	if *useDefaultCredentials {
		if *configuration != "" {
			log.Fatalf("-use-default-credentials and -configuration are mutual exclusive")
		}
		updater.credentials, err = google.FindDefaultCredentials(updater.ctx)
	} else {
		updater.credentials, err = gcloudconfig.GetCredentials(*configuration)
	}

	if project != "" {
		updater.credentials.ProjectID = project
	}

	if updater.credentials.ProjectID == "" {
		log.Fatalf("no project specified and no default project set")
	}
	if err != nil {
		log.Fatal(err)
	}

	err = updater.SourceImageDefinition()
	if err != nil {
		log.Fatal(err)
	}
}
