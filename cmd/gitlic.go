package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	ao "github.com/appoptics/appoptics-api-go"
	"github.com/solarwinds/gitlic-check/config"
	"github.com/solarwinds/gitlic-check/gitlic"
	"github.com/solarwinds/gitlic-check/swgithub"
	"github.com/spf13/cobra"
)

var (
	uploadOnly  bool
	noUpload    bool
	statsReport bool
)

func init() {
	gitlicCmd.Flags().BoolVar(&uploadOnly, "upload-only", false, "test upload only; skip GitHub check")
	gitlicCmd.Flags().BoolVar(&noUpload, "no-upload", false, "test/re-run GitHub check only; skip upload")
	gitlicCmd.Flags().BoolVar(&statsReport, "stats", false, "run a stats report instead of standard Gitlic")
	rootCmd.AddCommand(gitlicCmd)
}

var gitlicCmd = &cobra.Command{
	Use: "gitlic",
	Run: func(cmd *cobra.Command, args []string) {
		aoToken := os.Getenv("AO_TOKEN")
		aoClient := ao.NewClient(aoToken)
		mService := aoClient.MeasurementsService()
		measurement := ao.Measurement{
			Name:  "augit.gitlic.runs",
			Value: 1,
			Time:  time.Now().Unix(),
			Tags:  map[string]string{"environment": os.Getenv("ENVIRONMENT")},
		}
		run()
		log.Info(generateSuccessString("gitlic"))
		_, err := mService.Create(ao.NewMeasurementsBatch([]ao.Measurement{measurement}, nil))
		if err != nil {
			log.Fatalln(err)
		}
	},
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
	cf := config.GetConfig()
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory. Error: %v\n", err)
	}

	if uploadOnly && cf.Drive == nil {
		log.Fatalln("To use the test-upload flag, you must specify the relevant config parameters. See README.")
	}

	var filesToOutput []string
	if statsReport {
		filesToOutput = []string{"stats.csv"}
	} else {
		filesToOutput = []string{"repos.csv", "users.csv", "invites.csv"}
	}
	fo := prepareOutput(uploadOnly, filesToOutput)
	defer func() {
		for _, file := range fo {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	if !uploadOnly {
		if statsReport {
			swgithub.RunGitlicStatsReport(ctx, cf, fo)
		} else {
			swgithub.RunGitlicCheck(ctx, cf, fo)
		}
		for _, file := range fo {
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				log.Fatal(err)
			}
		}
	}

	if !noUpload && cf.Drive != nil {
		gitlic.UploadToDrive(ctx, cf, wd, fo)
	}
}
