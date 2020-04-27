package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Configuration struct {
		ActiveConfiguration string `json:"active_configuration"`
		Properties          struct {
			Builds struct {
				UseKaniko string `json:"use_kaniko"`
			} `json:"builds"`
			Compute struct {
				Region string `json:"region"`
				Zone   string `json:"zone"`
			} `json:"compute"`
			Core struct {
				Account               string `json:"account"`
				DisableUsageReporting string `json:"disable_usage_reporting"`
				Project               string `json:"project"`
			} `json:"core"`
			Run struct {
				Cluster         string `json:"cluster"`
				ClusterLocation string `json:"cluster_location"`
				Platform        string `json:"platform"`
				Region          string `json:"region"`
			} `json:"run"`
		} `json:"properties"`
	} `json:"configuration"`
	Credential struct {
		AccessToken string    `json:"access_token"`
		IDToken     string    `json:"id_token"`
		TokenExpiry time.Time `json:"token_expiry"`
	} `json:"credential"`
	Sentinels struct {
		ConfigSentinel string `json:"config_sentinel"`
	} `json:"sentinels"`
}

func getDefaultProject() string {
	var stdout, stderr bytes.Buffer
	result := os.Getenv("CLOUDSDK_CORE_PROJECT")
	exe, err := exec.LookPath("gcloud")
	if err != nil {
		log.Printf("gcloud not found, using environment variable CLOUDSDK_CORE_PROJECT")
		return result
	}

	cmd := exec.Command(exe, "config", "config-helper", "--format", "json")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Command %s failed with %s\n", cmd.String(), err)
		return result
	}

	var config Config
	err = json.Unmarshal(stdout.Bytes(), &config)
	if err != nil {
		log.Fatalf("failed to parse json output of %s into configuration, %s\n", cmd.String(), err)
	}

	if config.Configuration.Properties.Core.Project != "" {
		result = config.Configuration.Properties.Core.Project
	}
	log.Printf("default project is %s", result)
	return result
}

