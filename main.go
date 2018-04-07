package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type ghconfig struct {
	Token          string   `json:"pat"`
	IgnoredOrgs    []string `json:"ignoredOrgs,omitempty"`
	RmInvitesAfter int      `json:"rmInvitesAfter,omitempty"` // in hours
}

type driveconfig struct {
	OutputDir       string `json:"outputDir"`
	EnableTeamDrive bool   `json:"enableTeamDrive,omitempty"`
}

type config struct {
	Github *ghconfig    `json:"github,omitempty"`
	Drive  *driveconfig `json:"drive,omitempty"`
}

func getConfig() config {
	fh, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Failed to read config file. Error: %s\n", err)
	}
	var cf config
	if err := json.Unmarshal(fh, &cf); err != nil {
		log.Fatalf("Failed to parse config file. Error: %v\n", err)
	}
	// Ensure GitHub credentials have been included in config file
	if cf.Github == nil || cf.Github.Token == "" {
		log.Fatalf("Failed to parse PAT for GitHub. Please ensure you're following instructions in README.")
	}
	return cf
}

func prepareOutput(uploadOnly *bool, filesToOutput []string) map[string]*os.File {
	if *uploadOnly == false {
		os.RemoveAll("output")
		os.Mkdir("output", 0777)
	}
	os.Chdir("output")
	wd, _ := os.Getwd()
	fmt.Printf("WD: %s\n", wd)
	fo := make(map[string]*os.File)
	for _, name := range filesToOutput {
		file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Printf("Failed to create %s with %s\n", name, err)
		}
		fo[name] = file
	}
	os.Chdir("../")
	return fo
}

func main() {
	fUploadOnly := flag.Bool("upload-only", false, "test upload only; skip GitHub check")
	fNoUpload := flag.Bool("no-upload", false, "test/re-run GitHub check only; skip upload")
	flag.Parse()

	cf := getConfig()
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory. Error: %v\n", err)
	}

	if *fUploadOnly == true && cf.Drive == nil {
		log.Fatalln("To use the test-upload flag, you must specify the relevant config parameters. See README.\n")
	}
	filesToOutput := []string{"repos.csv", "users.csv", "invites.csv"}
	fo := prepareOutput(fUploadOnly, filesToOutput)
	defer func() {
		for _, file := range fo {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if *fUploadOnly == false {
		RunCheck(ctx, cf, fo)
		for _, file := range fo {
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				log.Fatal(err)
			}
		}

	}

	if *fNoUpload == false && cf.Drive != nil {
		UploadToDrive(ctx, cf, wd, fo)
	}
}
