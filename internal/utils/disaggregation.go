package utils

import (
	"CNGDataManager/internal/google/sheets"
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
