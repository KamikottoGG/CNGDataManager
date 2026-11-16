package utils

import (
	"CNGDataManager/internal/google/sheets/v2"
	"io"

	"encoding/json"
	"errors"
	"net/http"
	url2 "net/url"
	"strconv"

	"time"
)

type GetBeneficiariesDataParams struct {
	DateStart time.Time
	DateEnd   time.Time
}

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *APIClient) GetEvacuationRecords(params GetBeneficiariesDataParams) ([]sheets.EvacuationRecord, time.Time, error) {
	return fetchRecords[sheets.EvacuationRecord](c, "evacuation", params)
}

func (c *APIClient) GetTransitHubRecords(params GetBeneficiariesDataParams) ([]sheets.TransitHubRecord, time.Time, error) {
	return fetchRecords[sheets.TransitHubRecord](c, "transit-hub", params)
}

func (c *APIClient) GetLocalTransportationsRecords(params GetBeneficiariesDataParams) ([]sheets.LocalTransportationRecord, time.Time, error) {
	return fetchRecords[sheets.LocalTransportationRecord](c, "local-transportations", params)
}

func (c *APIClient) GetVehicleTripRecords(params GetBeneficiariesDataParams) ([]sheets.VehicleTripRecord, time.Time, error) {
	return fetchRecords[sheets.VehicleTripRecord](c, "vehicle-trip", params)
}

func (c *APIClient) GetMainShelterRecords(params GetBeneficiariesDataParams) ([]sheets.MainShelterRecord, time.Time, error) {
	return fetchRecords[sheets.MainShelterRecord](c, "shelter", params)
}

func (c *APIClient) GetHotlineRecords(params GetBeneficiariesDataParams) ([]sheets.HotlineCallRecord, time.Time, error) {
	return fetchRecords[sheets.HotlineCallRecord](c, "hotline", params)
}

type ReportData struct {
	StartDate            time.Time
	EndDate              time.Time
	Hotline              []sheets.HotlineCallRecord
	Evacuation           []sheets.EvacuationRecord
	Shelter              []sheets.MainShelterRecord
	AllShelter           []sheets.MainShelterRecord
	TransitHub           []sheets.TransitHubRecord
	AllTransitHub        []sheets.TransitHubRecord
	LocalTransportations []sheets.LocalTransportationRecord
	VehicleTrip          []sheets.VehicleTripRecord
	CacheUpdateTime      time.Time
}

func (c *APIClient) GetAllData(params GetBeneficiariesDataParams) (ReportData, error) {
	evac, _, err := c.GetEvacuationRecords(params)
	if err != nil {
		return ReportData{}, err
	}
	transit, _, err := c.GetTransitHubRecords(params)
	if err != nil {
		return ReportData{}, err
	}
	local, _, err := c.GetLocalTransportationsRecords(params)
	if err != nil {
		return ReportData{}, err
	}
	shelter, _, err := c.GetMainShelterRecords(params)
	if err != nil {
		return ReportData{}, err
	}
	hotline, _, err := c.GetHotlineRecords(params)
	if err != nil {
		return ReportData{}, err
	}
	vTrips, lastEditTime, err := c.GetVehicleTripRecords(params)
	if err != nil {
		return ReportData{}, err
	}

	allShelter, _, err := c.GetMainShelterRecords(GetBeneficiariesDataParams{})
	if err != nil {
		return ReportData{}, err
	}

	trasitHub, _, err := c.GetTransitHubRecords(GetBeneficiariesDataParams{})
	if err != nil {
		return ReportData{}, err
	}

	return ReportData{
		Evacuation:           evac,
		TransitHub:           transit,
		Shelter:              shelter,
		LocalTransportations: local,
		Hotline:              hotline,
		VehicleTrip:          vTrips,
		CacheUpdateTime:      lastEditTime,
		AllShelter:           allShelter,
		AllTransitHub:        trasitHub,
		StartDate:            params.DateStart,
		EndDate:              params.DateEnd,
	}, nil
}

func fetchRecords[T any](c *APIClient, endpoint string, params GetBeneficiariesDataParams) ([]T, time.Time, error) {
	urlStr, err := url2.JoinPath(c.baseURL, endpoint)

	if err != nil {
		return nil, time.Time{}, err
	}

	url, err := url2.Parse(urlStr)
	if err != nil {
		return nil, time.Time{}, err
	}
	q := url.Query()
	if !params.DateEnd.IsZero() {
		q.Add("endDate", params.DateEnd.Format("02-01-2006"))
	}

	if !params.DateStart.IsZero() {
		q.Add("startDate", params.DateStart.Format("02-01-2006"))
	}

	url.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(url.String())

	if err != nil {
		return nil, time.Time{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil, time.Time{}, errors.New("server return http code " + strconv.Itoa(resp.StatusCode))
	}

	records := []T{}
	w := json.NewDecoder(resp.Body)
	err = w.Decode(&records)

	if err != nil {
		return nil, time.Time{}, err
	}

	lastModified, err := time.Parse(time.DateTime, resp.Header.Get("Last-Modified"))

	if err != nil {
		return nil, time.Time{}, err
	}

	return records, lastModified, nil

}
