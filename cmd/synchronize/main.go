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

	eid := os.Getenv("EVACUATION_SPREADSHEET_ID")

	if eid == "" {
		log.Fatal("EVACUATION_SPREADSHEET_ID environment variable not set")
	}

	msid := os.Getenv("MAIN_SHELTER_SPREADSHEET_ID")

	if msid == "" {
		log.Fatal("MAIN_SHELTER_SPREADSHEET_ID environment variable not set")
	}

	thid := os.Getenv("TRANSIT_HUB_SPREADSHEET_ID")

	if thid == "" {
		log.Fatal("TRANSIT_HUB_SPREADSHEET_ID environment variable not set")
	}

	hid := os.Getenv("HOTLINE_SPREADSHEET_ID")

	if hid == "" {
		log.Fatal("HOTLINE_SPREADSHEET_ID environment variable not set")
	}
	//
	//eJson := os.Getenv("EVACUATION_CACHE_JSON_FILE_PATH")
	//
	//if eJson == "" {
	//	log.Fatal("EVACUATION_CACHE_JSON_FILE_PATH environment variable not set")
	//}

	mainShelterJson := os.Getenv("MAIN_SHELTER_CACHE_JSON_FILE_PATH")

	if mainShelterJson == "" {
		log.Fatal("MAIN_SHELTER_CACHE_JSON_FILE_PATH environment variable not set")
	}

	evacuationJsonPath := os.Getenv("EVACUATION_CACHE_JSON_FILE_PATH")

	if evacuationJsonPath == "" {
		log.Fatal("EVACUATION_CACHE_JSON_FILE_PATH environment variable not set")
	}

	transitHubJsonPath := os.Getenv("TRANSIT_HUB_CACHE_JSON_FILE_PATH")

	if transitHubJsonPath == "" {
		log.Fatal("TRANSIT_HUB_CACHE_JSON_FILE_PATH environment variable not set")
	}

	hotlineJsonPath := os.Getenv("HOTLINE_CACHE_JSON_FILE_PATH")

	if hotlineJsonPath == "" {
		log.Fatal("HOTLINE_CACHE_JSON_FILE_PATH environment variable not set")
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

	err = sheets.CacheEvacuationData(sheetsService, evacuationSpreadSheet, evacuationJsonPath)

	if err != nil {
		log.Fatal("[X]", "[Error while saving evacuation data]", err)
	}

	mainShelterSpreadSheet, err := sheetsService.Spreadsheets.Get(msid).Do()

	if err != nil {
		log.Fatal(err)
	}

	err = sheets.CacheMainShelterData(sheetsService, mainShelterSpreadSheet, mainShelterJson)

	if err != nil {
		log.Fatal("[X]", "[Error while saving main shelter data]", err)
	}

	transitHubSpreadSheet, err := sheetsService.Spreadsheets.Get(thid).Do()

	if err != nil {
		log.Fatal(err)
	}

	err = sheets.CacheTransitHubData(sheetsService, transitHubSpreadSheet, transitHubJsonPath)

	if err != nil {
		log.Fatal("[X]", "[Error while saving transit hub data]", err)
	}

	hotlineSpreadSheet, err := sheetsService.Spreadsheets.Get(hid).Do()

	if err != nil {
		log.Fatal(err)
	}

	err = sheets.CacheHotlineData(sheetsService, hotlineSpreadSheet, hotlineJsonPath)

	if err != nil {
		log.Fatal("[X]", "[Error while saving hotline data]", err)
	}
}
