package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"regexp"
	"sort"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type GoogleComputeBuilder struct {
	Type        string `json:"type"`
	Zone        string `json:"zone"`
	AccountFile string `json:"account_file"`
	Image       string `json:"source_image"`
	Family      string `json:"source_image_family"`
	ProjectId   string `json:"source_image_project_id"`
}

func (builder *GoogleComputeBuilder) updateGoogleSourceImage(ctx context.Context, credentials *google.Credentials) (bool, error) {
	updated := false

	if builder.AccountFile != "" {
		content, err := ioutil.ReadFile(builder.AccountFile)
		if err != nil {
			log.Fatal(err)
		}
		credentials, err = google.CredentialsFromJSON(ctx, content)
		if err != nil {
			return false, err
		}
	}

	computeService, err := compute.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		return false, err
	}

	// determine basename
	name := builder.Image
	r := regexp.MustCompile("^(.+)(-v.+)$")
	parts := r.FindStringSubmatch(name)
	if len(parts) > 0 {
		// String version suffix of name for search
		name = parts[1]
	}

	// determine project for builder
	project := builder.ProjectId
	if project == "" && builder.Family != "" {
		if p := PublicProjects.FindProjectForName(builder.Family); p != nil {
			project = p.Project
		}
	}
	if project == "" && name != "" {
		if p := PublicProjects.FindProjectForName(name); p != nil {
			project = p.Project
		}
	}
	if project == "" {
		project = credentials.ProjectID
	}

	// find latest image of the same name
	var list *compute.ImageList
	namePattern := ".*"
	familyPattern := ".*"

	if name != "" {
		namePattern = fmt.Sprintf("%s-v.*", name)
	}
	if builder.Family != "" {
		familyPattern = builder.Family
	}
	token := ""
	images := make([]*compute.Image, 0, 256)
	for {
		if token == "" {
			filter := fmt.Sprintf("(name eq '%s') (family eq '%s')", namePattern, familyPattern)
			list, err = computeService.Images.List(project).Filter(filter).Do()
		} else {
			list, err = computeService.Images.List(project).PageToken(token).Do()
		}
		if err != nil {
			return false, fmt.Errorf("failed to query images for name %s and family %s in project %s, %s", namePattern, familyPattern, project, err)
		}
		images = append(images, list.Items...)
		token = list.NextPageToken
		if token == "" {
			break
		}
	}

	if len(images) == 0 {
		return false, fmt.Errorf("no images found with name %s and family %s in project %s\n", namePattern, familyPattern,
			project)
	}

	sort.Sort(byCreationTimestamp(images))
	image := images[0]
	if builder.Family != "" && builder.Family != image.Family {
		return false, fmt.Errorf("image %s is of family %s, not %s", image.Name, image.Family, builder.Family)
	}

	if builder.Image != image.Name {
		log.Printf("updating image from '%s' to '%s'", builder.Image, image.Name)
		updated = true
		builder.Image = image.Name
	}
	if updated && builder.Family == "" {
		log.Printf("setting image family to '%s'", image.Family)
		builder.Family = image.Family
	}
	if updated && builder.ProjectId == "" {
		log.Printf("setting source image project to '%s'", project)
		builder.ProjectId = project
	}

	return updated, nil
}

type byCreationTimestamp []*compute.Image

func (s byCreationTimestamp) Len() int {
	return len(s)
}
func (s byCreationTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCreationTimestamp) Less(i, j int) bool {
	return s[i].CreationTimestamp > s[j].CreationTimestamp
}
