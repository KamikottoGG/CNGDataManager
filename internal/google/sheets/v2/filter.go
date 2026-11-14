package sheets

import "strings"

func ValuesFilter(keys ...string) RowFilter {
	return func(m map[string]string) bool {

		for _, k := range keys {
			if strings.TrimSpace(m[k]) == "" {
				return false
			}
		}

		return true
	}
}

func EvacuationRowFilter() RowFilter {
	return ValuesFilter("Дата [ДД/ММ/РРРР]", "Стать", "Прізвище", "Ім'я")
}

func TransitHubRowFilter() RowFilter {
	return ValuesFilter("Прізвище", "Ім'я", "Стать", "Статус проживання")
}

func LocalTransportationsRowFilter() RowFilter {
	return ValuesFilter("Прізвище", "Ім'я")
}

func MainShelterRowFilter() RowFilter {
	return ValuesFilter("Внес дані", "Ім'я", "Статус")
}

func VehicleTripRowFilter() RowFilter {
	return ValuesFilter("Водій/відповідальний за виїзд", "Марка машини – номери")
}

func HotlineRowFilter() RowFilter {
	return ValuesFilter("Відповідальний")
}
