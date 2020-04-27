package main

import (
	"context"
	"fmt"
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

func (builder *GoogleComputeBuilder) updateGoogleSourceImage() (bool, error) {
	updated := false

	ctx := context.Background()
	credentials := getCredentials(ctx, builder.AccountFile)
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
	if project == "" {
		if builder.Family != "" {
			if p := PublicProjects.FindProjectForName(builder.Family); p != nil {
				project = p.Project
			}
		} else if name != "" {
			if p := PublicProjects.FindProjectForName(name); p != nil {
				project = p.Project
			}
		} else {
			project = credentials.ProjectID
		}
	}

	// find latest image of the same name
	var list *compute.ImageList
	if name != "" {
		list, err = computeService.Images.List(project).Filter(fmt.Sprintf("name eq \"%s.*\"", name)).Do()
		if err != nil || len(list.Items) == 0 {
			return false, fmt.Errorf("no images found with name %s in project %s", name, project)
		}
	} else {
		list, err = computeService.Images.List(project).Filter(fmt.Sprintf("family eq ^%s$", builder.Family)).Do()
		if err != nil || len(list.Items) == 0 {
			return false, fmt.Errorf("no images found of family %s in project %s", builder.Family, project)
		}
	}

	sort.Sort(byCreationTimestamp(list.Items))
	image := list.Items[0]
	if builder.Family != "" && builder.Family != image.Family {
		return false, fmt.Errorf("image %s is of family %s, not %s", image.Name, image.Family, builder.Family)
	}

	if builder.Image != image.Name {
		log.Printf("updating image from %s to %s", builder.Image, image.Name)
		updated = true
		builder.Image = image.Name
	}
	if updated && builder.Family == "" {
		log.Printf("setting image family to  %s", image.Family)
		builder.Family = image.Family
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
