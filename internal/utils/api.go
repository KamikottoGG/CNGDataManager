package utils

import (
	"CNGDataManager/internal/google/sheets"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GetBeneficiariesDataParams struct {
	DateStart time.Time
	DateEnd   time.Time
}

func GetAllData(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.EvacuationRecord,
	[]sheets.TransitHubRecord,
	[]sheets.LocalTransportationRecord,
	[]sheets.MainShelterRecord,
	[]sheets.HotlineCallRecord,
	time.Time,
	error) {
	var err error

	e, err := GetEvacuationBeneficiariesRecords(apiEndpoint, params)

	if err != nil {
		return nil, nil, nil, nil, nil, time.Time{}, fmt.Errorf("error getting evacuation beneficiaries data: %w", err)
	}

	th, err := GetTransitHubRecords(apiEndpoint, params)

	if err != nil {
		return nil, nil, nil, nil, nil, time.Time{}, fmt.Errorf("error getting transit hub beneficiaries data: %w", err)
	}

	lt, err := GetLocalTransportationsRecords(apiEndpoint, params)

	if err != nil {
		return nil, nil, nil, nil, nil, time.Time{}, fmt.Errorf("error getting local transportations beneficiaries data: %w", err)
	}

	ms, cUpdatedTime, err := GetMainShelterRecords(apiEndpoint, params)
	if err != nil {
		return nil, nil, nil, nil, nil, time.Time{}, fmt.Errorf("error getting main shelter beneficiaries data: %w", err)
	}

	h, err := GetHotlineRecords(apiEndpoint, params)
	if err != nil {
		return nil, nil, nil, nil, nil, time.Time{}, fmt.Errorf("error getting hotline beneficiaries data: %w", err)
	}

	return e, th, lt, ms, h, cUpdatedTime, nil
}

func GetEvacuationBeneficiariesRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.EvacuationRecord, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/evacuation?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	records := []sheets.EvacuationRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetMainShelterRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.MainShelterRecord, time.Time, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/shelter?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, time.Time{}, err
	}

	defer resp.Body.Close()

	records := []sheets.MainShelterRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, time.Time{}, fmt.Errorf("error decoding response json: %w", err)
	}

	cUpdatedTime, err := time.Parse(time.DateTime, resp.Header.Get("Last-Modified"))

	if err != nil {
		return nil, time.Time{}, fmt.Errorf("error getting cache update time: %w", err)
	}

	return records, cUpdatedTime, nil
}

func GetTransitHubRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.TransitHubRecord, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/transit?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	records := []sheets.TransitHubRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetVehicleTripRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.VehicleTripRecord, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/vehicle-trip?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	records := []sheets.VehicleTripRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetLocalTransportationsRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.LocalTransportationRecord, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/local-transportations?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	records := []sheets.LocalTransportationRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetHotlineRecords(apiEndpoint string, params GetBeneficiariesDataParams) ([]sheets.HotlineCallRecord, error) {

	url := strings.Builder{}

	url.WriteString(apiEndpoint)
	url.WriteString("/hotline?")

	if !params.DateStart.IsZero() {
		url.WriteString("startDate=")
		url.WriteString(params.DateStart.Format("02-01-2006"))
		url.WriteString("&")
	}

	if !params.DateEnd.IsZero() {
		url.WriteString("endDate=")
		url.WriteString(params.DateEnd.Format("02-01-2006"))
	}

	resp, err := http.Get(url.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	records := []sheets.HotlineCallRecord{}

	w := json.NewDecoder(resp.Body)

	err = w.Decode(&records)

	if err != nil {
		return nil, err
	}

	return records, nil
}
