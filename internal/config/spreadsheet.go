package config

import "errors"

const (
	ServiceAccountPath                = "SERVICE_ACCOUNT_FILE_PATH"
	EvacuationSpreadsheetID           = "EVACUATION_SPREADSHEET_ID"
	MainShelterSpreadsheetID          = "MAIN_SHELTER_SPREADSHEET_ID"
	TransitHubSpreadsheetID           = "TRANSIT_HUB_SPREADSHEET_ID"
	HotlineSpreadsheetID              = "HOTLINE_SPREADSHEET_ID"
	LocalTransportationsSpreadsheetID = "LOCAL_TRANSPORTATION_SPREADSHEET_ID"
	VehicleTripSpreadsheetID          = "VEHICLE_TRIP_SPREADSHEET_ID"
)

type SpreadsheetConfig struct {
	ServiceAccountPath string
	EvacuationID       string
	MainShelterID      string
	TransitHubID       string
	HotlineID          string
	LocalTransportID   string
	VehicleTripID      string
}

func NewSpreadsheetConfig(c Config) (SpreadsheetConfig, error) {
	var err error
	req := func(param string) string {
		if err != nil {
			return ""
		}

		v := c.Get(param)

		if v == "" {
			err = errors.New(param + " config required")
		}

		return v
	}

	ssc := SpreadsheetConfig{
		ServiceAccountPath: req(ServiceAccountPath),
		EvacuationID:       req(EvacuationSpreadsheetID),
		MainShelterID:      req(MainShelterSpreadsheetID),
		TransitHubID:       req(TransitHubSpreadsheetID),
		HotlineID:          req(HotlineSpreadsheetID),
		LocalTransportID:   req(LocalTransportationsSpreadsheetID),
		VehicleTripID:      req(VehicleTripSpreadsheetID),
	}

	if err != nil {
		return SpreadsheetConfig{}, err
	}

	return ssc, nil
}
