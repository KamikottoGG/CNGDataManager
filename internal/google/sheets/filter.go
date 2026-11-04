package sheets

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func text(msg ...interface{}) string {
	return fmt.Sprintf("%+v", msg...)
}

func ExtractValuesFromRows(rows *sheets.ValueRange) *[]map[string]string {

	var values = make([]map[string]string, 0, len(rows.Values))

	columnHeaders := rows.Values[0]

	for _, row := range rows.Values[1:] {

		mappedObj := make(map[string]string)

		for i, value := range row {
			mappedObj[text(columnHeaders[i])] = text(value)
		}
		values = append(values, mappedObj)
	}

	return &values
}

func FilterValues(values *[]map[string]string) *[]map[string]string {
	var newValues = make([]map[string]string, 0, len(*values)/4)

	for _, value := range *values {
		if len(value["Дата [ДД/ММ/РРРР]"]) == 0 {
			continue
		}
		if len(value["Стать"]) == 0 {
			continue
		}
		if len(value["Прізвище"]) == 0 {
			continue
		}
		if len(value["Ім'я"]) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return &newValues

}
