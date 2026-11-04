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

	opts := []option.ClientOption{
		option.WithAuthCredentialsFile(
			"service_account", credFilePath,
		),
	}

	sheetsService, err := gsheets.NewService(context.Background(), opts...)

	if err != nil {
		log.Fatal(err)
	}

	eid := os.Getenv("EVACUATION_SPREAD_SHEET_ID")

	if eid == "" {
		log.Fatal("EVACUATION_SPREAD_SHEET_ID environment variable not set")
	}

	evacuationSpreadSheet, err := sheetsService.Spreadsheets.Get(eid).Do()

	if err != nil {
		log.Fatal(err)
	}

	log.Print("[+] ", "Evacuation spreadsheet got with name: "+evacuationSpreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(evacuationSpreadSheet.Sheets))

	collectedData := make([]*[]sheets.Evacuation, 0, len(evacuationSpreadSheet.Sheets))

	rowsAffected := 0

	for i, sheet := range evacuationSpreadSheet.Sheets {
		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(evacuationSpreadSheet.Sheets), sheet.Properties.Title)

		rows, err := sheetsService.Spreadsheets.Values.Get(evacuationSpreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("	[?] %d rows parsed. ", len(rows.Values))

		extractedValues := sheets.ExtractValuesFromRows(rows)

		filteredValues := sheets.FilterValues(extractedValues)

		fvr := len(filteredValues)

		log.Printf("	[?] %d rows filtered. ", fvr)

		mappedObjects := make([]sheets.Evacuation, 0, fvr)

		for j, v := range filteredValues {

			var mappedObject sheets.Evacuation

			if err := mappedObject.MapObject(v); err != nil {
				log.Printf("[!] row %d: %v", j, err)
				continue
			}

			mappedObjects = append(mappedObjects, mappedObject)

			rowsAffected++
		}

		collectedData = append(collectedData, &mappedObjects)

		log.Printf("[%d/%d] Work with sheet end. ", i+1, len(evacuationSpreadSheet.Sheets))

	}

	log.Printf("	[+] %d month iterated. ", len(collectedData))
	log.Printf("[+] %d rows mapped (not persisted yet)", rowsAffected)

}
