package main

import (
	"CNGDataManager/internal/config"
	"CNGDataManager/internal/google/sheets/v2"
	"context"
	"log"

	"google.golang.org/api/option"
	gsheets "google.golang.org/api/sheets/v4"
)

type SyncTask func() error

func NewSyncTask[T any](
	getter sheets.RowGetter,
	mapper sheets.RecordMapper[T],
	cacher sheets.Cacher[T],
	filters ...sheets.RowFilter) SyncTask {

	return func() error {
		records, err := sheets.GetAndMapRows(getter, mapper, filters...)
		if err != nil {
			return err
		}
		return cacher.Save(records)
	}

}

func main() {

	c := config.NewEnvConfig()

	spreadsheetConfig, err := config.NewSpreadsheetConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	cfc, err := config.NewCacheFilesConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	opts := []option.ClientOption{
		option.WithAuthCredentialsFile(
			option.ServiceAccount, spreadsheetConfig.ServiceAccountPath,
		),
	}

	s, err := gsheets.NewService(context.Background(), opts...)

	if err != nil {
		log.Fatal(err)
	}

	syncTasks := []SyncTask{
		NewSyncTask(sheets.NewEvacuationRowGetter(s, spreadsheetConfig.EvacuationID), sheets.EvacuationMapper, sheets.NewJSONCacher[sheets.EvacuationRecord](cfc.EvacuationJsonPath), sheets.EvacuationRowFilter()),
		NewSyncTask(sheets.NewTransitHubRowGetter(s, spreadsheetConfig.TransitHubID), sheets.TransitHubMapper, sheets.NewJSONCacher[sheets.TransitHubRecord](cfc.TransitHubJsonPath), sheets.TransitHubRowFilter()),
		NewSyncTask(sheets.NewMainShelterRowGetter(s, spreadsheetConfig.MainShelterID), sheets.MainShelterMapper, sheets.NewJSONCacher[sheets.MainShelterRecord](cfc.MainShelterJsonPath), sheets.MainShelterRowFilter()),
		NewSyncTask(sheets.NewLocalTransportationsRowGetter(s, spreadsheetConfig.LocalTransportID), sheets.LocalTransportationsMapper, sheets.NewJSONCacher[sheets.LocalTransportationRecord](cfc.LocalTransportationsJsonPath), sheets.LocalTransportationsRowFilter()),
		NewSyncTask(sheets.NewVehicleTripRowGetter(s, spreadsheetConfig.VehicleTripID), sheets.VehicleTripMapper, sheets.NewJSONCacher[sheets.VehicleTripRecord](cfc.VehicleTripJsonPath), sheets.VehicleTripRowFilter()),
		NewSyncTask(sheets.NewHotlineRowGetter(s, spreadsheetConfig.HotlineID), sheets.HotlineMapper, sheets.NewJSONCacher[sheets.HotlineCallRecord](cfc.HotlineJsonPath), sheets.HotlineRowFilter()),
	}

	for _, syncTask := range syncTasks {
		if err := syncTask(); err != nil {
			log.Fatal(err)
		}
	}

}

//
//package main
//
//import (
//"CNGDataManager/internal/config"
//"CNGDataManager/internal/google/sheets/v2"
//"context"
//"log"
//"os"
//"strings"
//
//"google.golang.org/api/option"
//gsheets "google.golang.org/api/sheets/v4"
//)
//
//func main() {
//
//	c := config.NewEnvConfig()
//
//	spreadsheetConfig, err := config.NewSpreadsheetConfig(c)
//	if err != nil {
//		log.Fatal(err)
//	}
//	cacheFilesConfig, err := config.NewCacheFilesConfig(c)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	opts := []option.ClientOption{
//		option.WithAuthCredentialsFile(
//			"service_account", spreadsheetConfig.ServiceAccountPath,
//		),
//	}
//
//	sheetsService, err := gsheets.NewService(context.Background(), opts...)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	getSpreadsheet := func (id string)(*gsheets.Spreadsheet, error){
//		return sheetsService.Spreadsheets.Get(id).Do()
//	}
//
//	evacuationSpreadSheet, err := getSpreadsheet()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheEvacuationData(sheetsService, evacuationSpreadSheet, evacuationJsonPath)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving evacuation data]", err)
//	}
//
//	mainShelterSpreadSheet, err := sheetsService.Spreadsheets.Get(msid).Do()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheMainShelterData(sheetsService, mainShelterSpreadSheet, mainShelterJson)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving main shelter data]", err)
//	}
//
//	transitHubSpreadSheet, err := sheetsService.Spreadsheets.Get(thid).Do()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheTransitHubData(sheetsService, transitHubSpreadSheet, transitHubJsonPath)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving transit hub data]", err)
//	}
//
//	hotlineSpreadSheet, err := sheetsService.Spreadsheets.Get(hid).Do()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheHotlineData(sheetsService, hotlineSpreadSheet, hotlineJsonPath)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving hotline data]", err)
//	}
//
//	lTransportationsSpreadSheet, err := sheetsService.Spreadsheets.Get(ltid).Do()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheLocalTransportationData(sheetsService, lTransportationsSpreadSheet, localTransportationsJsonPath)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving hotline data]", err)
//	}
//
//	vehicleTripSpreadsheet, err := sheetsService.Spreadsheets.Get(vtid).Do()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = sheets.CacheVehicleTripData(sheetsService, vehicleTripSpreadsheet, vehicleTripJsonPath)
//
//	if err != nil {
//		log.Fatal("[X]", "[Error while saving vehicle trips data]", err)
//	}
//
//	// CacheVehicleTripData
//
//	pathes := []string{
//		evacuationJsonPath,
//		transitHubJsonPath,
//		hotlineJsonPath,
//		localTransportationsJsonPath,
//		mainShelterJson,
//		vehicleTripJsonPath,
//	}
//
//	for _, p := range pathes {
//
//		n := strings.ReplaceAll(p, ".tmp", "")
//
//		err = os.Rename(p, n)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}
