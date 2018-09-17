package gitlic

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/solarwinds/gitlic-check/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

// UploadToDrive initiates the process of uploading your files to the specified Drive directory.
func UploadToDrive(ctx context.Context, cf config.Config, wd string, fo map[string]*os.File) {
	log.Print("Uploading to Google Drive...\n")

	secret, err := ioutil.ReadFile(filepath.Join(wd, "config/drive-key.json"))
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

func UploadToSheets(values [][]interface{}, dcf *config.DriveConfig) error {
	wd, _ := os.Getwd()
	secret, err := ioutil.ReadFile(filepath.Join(wd, "config/drive-key.json"))
	if err != nil {
		log.Fatalf("Failed to read JSON config file. Error: %v\n", err)
	}
	config, err := google.JWTConfigFromJSON(secret, drive.DriveScope)
	if err != nil {
		log.Fatalf("Failed to generate JWT from config file. Error: %v\n", err)
	}
	sheetClient, _ := sheets.New(config.Client(context.Background()))
	clear := &sheets.BatchClearValuesRequest{
		Ranges: []string{"A2:B"},
	}
	clearCall := sheetClient.Spreadsheets.Values.BatchClear(dcf.GhSheetId, clear)
	_, err = clearCall.Do()
	if err != nil {
		return err
	}
	vr := []*sheets.ValueRange{
		{
			Range:  "A2:B",
			Values: values,
		},
	}
	batch := &sheets.BatchUpdateValuesRequest{
		Data:             vr,
		ValueInputOption: "USER_ENTERED",
	}
	updateCall := sheetClient.Spreadsheets.Values.BatchUpdate(dcf.GhSheetId, batch)
	_, err = updateCall.Do()
	if err != nil {
		return err
	}
	resize := &sheets.AutoResizeDimensionsRequest{
		Dimensions: &sheets.DimensionRange{
			Dimension:  "COLUMNS",
			StartIndex: 0,
			EndIndex:   2,
		},
	}
	resizeReq := &sheets.Request{
		AutoResizeDimensions: resize,
	}
	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{resizeReq},
	}
	resizeCall := sheetClient.Spreadsheets.BatchUpdate(dcf.GhSheetId, req)
	_, err = resizeCall.Do()
	return err
}
