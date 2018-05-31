package gitlic

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

// UploadToDrive initiates the process of uploading your files to the specified Drive directory.
func UploadToDrive(ctx context.Context, cf Config, wd string, fo map[string]*os.File) {
	log.Print("Uploading to Google Drive...\n")

	secret, err := ioutil.ReadFile(filepath.Join(wd, "config/config-drive.json"))
	if err != nil {
		log.Fatalf("Failed to read JSON config file. Error: %v\n", err)
	}
	config, err := google.JWTConfigFromJSON(secret, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Failed to generate JWT from config file. Error: %v\n", err)
	}
	dcl, err := drive.New(config.Client(ctx))
	if err != nil {
		log.Fatalf("Failed to connect to Drive: %v\n", err)
	}

	parents := []string{cf.Drive.OutputDir}
	for name, file := range fo {
		f := &drive.File{
			MimeType: "text/csv",
			Name:     name,
			Parents:  parents,
		}
		_, err := dcl.Files.Create(f).Media(file).SupportsTeamDrives(cf.Drive.EnableTeamDrive).Do()
		if err != nil {
			log.Printf("Failed to upload %s: %s\n", file.Name(), err)
		} else {
			log.Printf("Successfully uploaded %s", file.Name())
		}
	}
	return
}
