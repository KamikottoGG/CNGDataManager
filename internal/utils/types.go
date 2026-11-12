package utils

type GroupCount struct {
	Total          int `json:"total"`
	WithDisability int `json:"with_disability"`
}

func (g *GroupCount) Add(disability bool) {
	g.Total++
	if disability {
		g.WithDisability++
	}
}

type BeneficiariesDisaggregation struct {
	Boys         GroupCount `json:"boys"`
	Girls        GroupCount `json:"girls"`
	Men          GroupCount `json:"men"`
	Women        GroupCount `json:"women"`
	ElderlyMen   GroupCount `json:"elderly_men"`
	ElderlyWomen GroupCount `json:"elderly_women"`
}

func (b *BeneficiariesDisaggregation) Total() int {
	return b.Boys.Total + b.Girls.Total + b.Men.Total + b.Women.Total + b.ElderlyMen.Total + b.ElderlyWomen.Total
}

//func IsEvacuationBeneficiariesCountSame(vehicleTripRecords []sheets.VehicleTripRecord, evacuationRecords []sheets.EvacuationRecord) bool{
//
//	vehicleEvacuatedCount := 0
//
//	for _, r := range vehicleTripRecords {
//		vehicleEvacuatedCount += r.EvacueesCount
//	}
//
//	return true
//}
