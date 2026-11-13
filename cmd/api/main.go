package main

import (
	"CNGDataManager/internal/google/sheets"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

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

	localTransportationsJsonPath := os.Getenv("LOCAL_TRANSPORTATION_CACHE_JSON_FILE_PATH")

	if localTransportationsJsonPath == "" {
		log.Fatal("LOCAL_TRANSPORTATION_CACHE_JSON_FILE_PATH environment variable not set")
	}

	vehicleTripJsonPath := os.Getenv("VEHICLE_TRIP_CACHE_JSON_FILE")

	if vehicleTripJsonPath == "" {
		log.Fatal("VEHICLE_TRIP_CACHE_JSON_FILE environment variable not set")
	}

	http.HandleFunc("/evacuation", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(evacuationJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(evacuationJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var evacuationData []sheets.EvacuationRecord
		err = json.Unmarshal(data, &evacuationData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var responseEvacuationData = make([]sheets.EvacuationRecord, 0, len(evacuationData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, evacuationRecord := range evacuationData {

			activityDate := time.Time(evacuationRecord.Date)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			responseEvacuationData = append(responseEvacuationData, evacuationRecord)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(responseEvacuationData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/hotline", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(hotlineJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(hotlineJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var hData []sheets.HotlineCallRecord
		err = json.Unmarshal(data, &hData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var respHotlineData = make([]sheets.HotlineCallRecord, 0, len(hData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, evacuationRecord := range hData {

			activityDate := time.Time(evacuationRecord.CallDate)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			respHotlineData = append(respHotlineData, evacuationRecord)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(respHotlineData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/local-transportations", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(localTransportationsJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(localTransportationsJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var tData []sheets.LocalTransportationRecord
		err = json.Unmarshal(data, &tData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var respTransData = make([]sheets.LocalTransportationRecord, 0, len(tData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, evacuationRecord := range tData {

			activityDate := time.Time(evacuationRecord.Date)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			respTransData = append(respTransData, evacuationRecord)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(respTransData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/shelter", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(mainShelterJson)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var sData []sheets.MainShelterRecord
		err = json.Unmarshal(data, &sData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(mainShelterJson)

		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var respShelterData = make([]sheets.MainShelterRecord, 0, len(sData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, evacuationRecord := range sData {

			activityDate := time.Time(evacuationRecord.ArrivalDate)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			respShelterData = append(respShelterData, evacuationRecord)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(respShelterData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/transit-hub", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(transitHubJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(transitHubJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var tData []sheets.TransitHubRecord
		err = json.Unmarshal(data, &tData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var respTransitData = make([]sheets.TransitHubRecord, 0, len(tData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, evacuationRecord := range tData {

			activityDate := time.Time(evacuationRecord.ArrivalDate)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			respTransitData = append(respTransitData, evacuationRecord)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(respTransitData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/vehicle-trip", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()

		startDate := query.Get("startDate")
		endDate := query.Get("endDate")

		data, err := os.ReadFile(vehicleTripJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		info, err := os.Stat(vehicleTripJsonPath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var lastEdit = info.ModTime()

		var tData []sheets.VehicleTripRecord
		err = json.Unmarshal(data, &tData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var respData = make([]sheets.VehicleTripRecord, 0, len(tData)/8)

		startTime, err := time.Parse("02-01-2006", startDate)

		if err != nil {
			startTime = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
		}

		endTime, err := time.Parse("02-01-2006", endDate)
		if err != nil {
			endTime = time.Now()
		} else {

			endTime = endTime.Add(24*time.Hour - time.Nanosecond)
		}

		for _, record := range tData {

			activityDate := time.Time(record.Date)

			if activityDate.Before(startTime) || activityDate.After(endTime) {
				continue
			}

			respData = append(respData, record)

		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Last-Modified", lastEdit.Format(time.DateTime))
		if err := json.NewEncoder(w).Encode(respData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
