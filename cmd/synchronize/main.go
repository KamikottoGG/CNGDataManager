package main

import (
	"CNGDataManager/internal/google/sheets"
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"
)

func main() {

	credFilePath := os.Getenv("SERVICE_ACCOUNT_FILE_PATH")

	if credFilePath == "" {
		log.Fatal("SERVICE_ACCOUNT_FILE_PATH environment variable not set")
	}

	eid := os.Getenv("EVACUATION_SPREAD_SHEET_ID")

	if eid == "" {
		log.Fatal("EVACUATION_SPREAD_SHEET_ID environment variable not set")
	}

	eJson := os.Getenv("EVACUATION_CACHE_JSON_FILE_PATH")

	if eJson == "" {
		log.Fatal("EVACUATION_CACHE_JSON_FILE_PATH environment variable not set")
	}

	opts := []option.ClientOption{
		option.WithAuthCredentialsFile(
			"service_account", credFilePath,
		),
	}

	sheetsService, err := gsheets.NewService(context.Background(), opts...)

	if err != nil {
		log.Fatal(err)
	}

	evacuationSpreadSheet, err := sheetsService.Spreadsheets.Get(eid).Do()

	if err != nil {
		log.Fatal(err)
	}

	err = sheets.CacheEvacuationData(sheetsService, evacuationSpreadSheet, eJson)

	if err != nil {
		log.Fatal("[X]", "[Error while saving evacuation data]", err)
	}
}
