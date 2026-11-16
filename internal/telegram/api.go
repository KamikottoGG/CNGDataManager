package telegram

import (
	"CNGDataManager/internal/google/sheets/v2"
	"CNGDataManager/internal/utils"

	"time"
)

type BeneficiaryAPI interface {
	//GetEvacuationRecords(params utils.GetBeneficiariesDataParams) ([]sheets.EvacuationRecord, time.Time, error)
	GetTransitHubRecords(params utils.GetBeneficiariesDataParams) ([]sheets.TransitHubRecord, time.Time, error)
	//GetLocalTransportationsRecords(params utils.GetBeneficiariesDataParams) ([]sheets.LocalTransportationRecord, time.Time, error)
	//GetVehicleTripRecords(params utils.GetBeneficiariesDataParams) ([]sheets.VehicleTripRecord, time.Time, error)
	GetMainShelterRecords(params utils.GetBeneficiariesDataParams) ([]sheets.MainShelterRecord, time.Time, error)
	//GetHotlineRecords(params utils.GetBeneficiariesDataParams) ([]sheets.HotlineCallRecord, time.Time, error)
	GetAllData(params utils.GetBeneficiariesDataParams) (utils.ReportData, error)
}
