package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Github *GhConfig    `json:"github,omitempty"`
	Drive  *DriveConfig `json:"drive,omitempty"`
}

type GhConfig struct {
	Token          string   `json:"pat"`
	IgnoredOrgs    []string `json:"ignoredOrgs,omitempty"`
	RmInvitesAfter int      `json:"rmInvitesAfter,omitempty"` // in hours
}

type DriveConfig struct {
	OutputDir       string `json:"outputDir"`
	EnableTeamDrive bool   `json:"enableTeamDrive,omitempty"`
	GhSheetId       string `json:"ghSheetId"`
}

func GetConfig() Config {
	fh, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Failed to read config file. Error: %s\n", err)
	}
	var cf Config
	if err := json.Unmarshal(fh, &cf); err != nil {
		log.Fatalf("Failed to parse config file. Error: %v\n", err)
	}
	// Ensure GitHub credentials have been included in config file
	if cf.Github == nil || cf.Github.Token == "" {
		log.Fatalf("Failed to parse PAT for GitHub. Please ensure you're following instructions in README.")
	}
	return cf
}
