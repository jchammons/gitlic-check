package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"github.com/solarwinds/gitlic-check/gitlic"
	"github.com/spf13/cobra"
)

var (
	uploadOnly bool
	noUpload bool
)

func init() {
	gitlicCmd.Flags().BoolVar(&uploadOnly, "upload-only", false, "test upload only; skip GitHub check")
	gitlicCmd.Flags().BoolVar(&noUpload, "no-upload", false, "test/re-run GitHub check only; skip upload")
	rootCmd.AddCommand(gitlicCmd)
}

var gitlicCmd = &cobra.Command{
	Use: "gitlic",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func getConfig() gitlic.Config {
	fh, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("Failed to read config file. Error: %s\n", err)
	}
	var cf gitlic.Config
	if err := json.Unmarshal(fh, &cf); err != nil {
		log.Fatalf("Failed to parse config file. Error: %v\n", err)
	}
	// Ensure GitHub credentials have been included in config file
	if cf.Github == nil || cf.Github.Token == "" {
		log.Fatalf("Failed to parse PAT for GitHub. Please ensure you're following instructions in README.")
	}
	return cf
}

func prepareOutput(uploadOnly bool, filesToOutput []string) map[string]*os.File {
	if uploadOnly == false {
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

func run() {
	cf := getConfig()
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory. Error: %v\n", err)
	}

	if uploadOnly == true && cf.Drive == nil {
		log.Fatalln("To use the test-upload flag, you must specify the relevant config parameters. See README.\n")
	}
	filesToOutput := []string{"repos.csv", "users.csv", "invites.csv"}
	fo := prepareOutput(uploadOnly, filesToOutput)
	defer func() {
		for _, file := range fo {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if uploadOnly == false {
		gitlic.RunCheck(ctx, cf, fo)
		for _, file := range fo {
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				log.Fatal(err)
			}
		}

	}

	if uploadOnly == false && cf.Drive != nil {
		gitlic.UploadToDrive(ctx, cf, wd, fo)
	}
}
