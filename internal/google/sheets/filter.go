package sheets

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func text(msg ...interface{}) string {
	return fmt.Sprintf("%+v", msg...)
}

func ExtractValuesFromRows(rows *sheets.ValueRange) []map[string]string {

	var values = make([]map[string]string, 0, len(rows.Values))

	columnHeaders := rows.Values[0]

	for _, row := range rows.Values[1:] {

		mappedObj := make(map[string]string)

		for i, value := range row {
			mappedObj[text(columnHeaders[i])] = text(value)
		}
		values = append(values, mappedObj)
	}

	return values
}
