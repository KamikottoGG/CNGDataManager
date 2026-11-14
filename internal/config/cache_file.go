package config

import "errors"

const (
	MainShelterJson          = "MAIN_SHELTER_CACHE_JSON_FILE_PATH"
	EvacuationJson           = "EVACUATION_CACHE_JSON_FILE_PATH"
	TransitHubJson           = "TRANSIT_HUB_CACHE_JSON_FILE_PATH"
	HotlineJson              = "HOTLINE_CACHE_JSON_FILE_PATH"
	LocalTransportationsJson = "LOCAL_TRANSPORTATION_CACHE_JSON_FILE_PATH"
	VehicleTripJson          = "VEHICLE_TRIP_CACHE_JSON_FILE"
)

type CacheFilesConfig struct {
	MainShelterJsonPath          string
	EvacuationJsonPath           string
	TransitHubJsonPath           string
	HotlineJsonPath              string
	LocalTransportationsJsonPath string
	VehicleTripJsonPath          string
}

func NewCacheFilesConfig(c Config) (CacheFilesConfig, error) {
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

	cfc := CacheFilesConfig{
		MainShelterJsonPath:          req(MainShelterJson),
		EvacuationJsonPath:           req(EvacuationJson),
		TransitHubJsonPath:           req(TransitHubJson),
		HotlineJsonPath:              req(HotlineJson),
		LocalTransportationsJsonPath: req(LocalTransportationsJson),
		VehicleTripJsonPath:          req(VehicleTripJson),
	}

	if err != nil {
		return CacheFilesConfig{}, err
	}

	return cfc, nil
}
