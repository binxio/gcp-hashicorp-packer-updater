package main

import (
	"context"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
)

func getCredentials(ctx context.Context, accountFile string) *google.Credentials {

	var credentials *google.Credentials
	var err error
	if accountFile != "" {
		content, err := ioutil.ReadFile(accountFile)
		if err != nil {
			log.Fatal(err)
		}
		credentials, err = google.CredentialsFromJSON(ctx, content)
	} else {
		credentials, err = google.FindDefaultCredentials(ctx)
	}
	if err != nil {
		log.Fatal(err)
	}

	if credentials.ProjectID == "" {
		credentials.ProjectID = getDefaultProject()
	}

	return credentials
}

