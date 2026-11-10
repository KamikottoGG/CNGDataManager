package utils

import (
	"CNGDataManager/internal/google/sheets"
	"strings"
	"unicode"
	"unicode/utf8"
)

func GetEvacuationBeneficiariesDisaggregation(records []sheets.EvacuationRecord) BeneficiariesDisaggregation {

	disaggregation := BeneficiariesDisaggregation{}

	for _, r := range records {

		switch {
		case r.Age < 18:
			if r.Sex == sheets.MALE {
				disaggregation.Boys.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Girls.Add(r.Disability)
			}
		case r.Age >= 18 && r.Age < 60:
			if r.Sex == sheets.MALE {
				disaggregation.Men.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Women.Add(r.Disability)
			}

		default:
			if r.Sex == sheets.MALE {
				disaggregation.ElderlyMen.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.ElderlyWomen.Add(r.Disability)
			}

		}

	}

	return disaggregation
}

func GetTransitHubBeneficiariesDisaggregation(records []sheets.TransitHubRecord) BeneficiariesDisaggregation {

	disaggregation := BeneficiariesDisaggregation{}

	for _, r := range records {

		switch {
		case r.Age < 18:
			if r.Sex == sheets.MALE {
				disaggregation.Boys.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Girls.Add(r.Disability)
			}
		case r.Age >= 18 && r.Age < 60:
			if r.Sex == sheets.MALE {
				disaggregation.Men.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Women.Add(r.Disability)
			}

		default:
			if r.Sex == sheets.MALE {
				disaggregation.ElderlyMen.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.ElderlyWomen.Add(r.Disability)
			}

		}

	}

	return disaggregation
}

func GetMainShelterBeneficiariesDisaggregation(records []sheets.MainShelterRecord) BeneficiariesDisaggregation {
	disaggregation := BeneficiariesDisaggregation{}

	for _, r := range records {

		switch {
		case r.Age < 18:
			if r.Sex == sheets.MALE {
				disaggregation.Boys.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Girls.Add(r.Disability)
			}
		case r.Age >= 18 && r.Age < 60:
			if r.Sex == sheets.MALE {
				disaggregation.Men.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Women.Add(r.Disability)
			}

		default:
			if r.Sex == sheets.MALE {
				disaggregation.ElderlyMen.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.ElderlyWomen.Add(r.Disability)
			}

		}

	}

	return disaggregation
}

func GetTransportationsBeneficiariesDisaggregation(records []sheets.LocalTransportationRecord) BeneficiariesDisaggregation {
	disaggregation := BeneficiariesDisaggregation{}

	for _, r := range records {

		switch {
		case r.Age < 18:
			if r.Sex == sheets.MALE {
				disaggregation.Boys.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Girls.Add(r.Disability)
			}
		case r.Age >= 18 && r.Age < 60:
			if r.Sex == sheets.MALE {
				disaggregation.Men.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.Women.Add(r.Disability)
			}

		default:
			if r.Sex == sheets.MALE {
				disaggregation.ElderlyMen.Add(r.Disability)
			}

			if r.Sex == sheets.FEMALE {
				disaggregation.ElderlyWomen.Add(r.Disability)
			}

		}

	}

	return disaggregation
}

func toTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {

		return "Невідомо"
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + strings.ToLower(s[size:])
}

func GetEvacuationGeoDisaggregation(records []sheets.EvacuationRecord) map[string]map[string]int {
	disaggregation := make(map[string]map[string]int)

	for _, r := range records {
		region := toTitle(string(r.Region))
		district := toTitle(string(r.District))

		if region == "Невідомо" {
			continue
		}

		if disaggregation[region] == nil {
			disaggregation[region] = make(map[string]int)
		}

		disaggregation[region][district]++
	}

	return (disaggregation)
}

func GetTransitHubGeoDisaggregation(records []sheets.TransitHubRecord) map[string]map[string]int {
	disaggregation := make(map[string]map[string]int)

	for _, r := range records {
		region := toTitle(string(r.FromRegion))
		district := toTitle(string(r.FromDistrict))

		if region == "Невідомо" {
			continue
		}

		if disaggregation[region] == nil {
			disaggregation[region] = make(map[string]int)
		}

		disaggregation[region][district]++
	}

	return (disaggregation)
}
