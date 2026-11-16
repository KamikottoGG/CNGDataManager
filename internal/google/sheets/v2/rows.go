package sheets

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func NewEvacuationRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("evacuation batch get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		for _, s := range spr.Sheets {
			f = append(f, &sheets.DataFilter{
				A1Range: fmt.Sprintf("'%s'!A1:W", s.Properties.Title),
			})
		}

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("evacuation batch get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func NewTransitHubRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("transit hub batch get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		for _, s := range spr.Sheets {
			f = append(f, &sheets.DataFilter{
				A1Range: fmt.Sprintf("'%s'!A1:W", s.Properties.Title),
			})
		}

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("transit hub batch get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func NewMainShelterRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("main shelter get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		f = append(f, &sheets.DataFilter{
			A1Range: "'Form_Responses'!A1:Z",
		})

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("main shelter get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func NewHotlineRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("hotline batch get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		for _, s := range spr.Sheets {
			f = append(f, &sheets.DataFilter{
				A1Range: fmt.Sprintf("'%s'!A1:W", s.Properties.Title),
			})
		}

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("hotline batch get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func NewVehicleTripRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("vehicle trip batch get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		for _, s := range spr.Sheets {
			if s.Properties.Title == "Відповіді форми (2)" {
				continue
			}

			f = append(f, &sheets.DataFilter{
				A1Range: fmt.Sprintf("'%s'!A1:R", s.Properties.Title),
			})
		}

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("vehicle trip batch get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func NewLocalTransportationsRowGetter(service *sheets.Service, spreadsheetID string) RowGetter {
	return func() ([]map[string]string, error) {
		spr, err := service.Spreadsheets.Get(spreadsheetID).Fields("spreadsheetId,sheets.properties.title").Do()
		if err != nil {
			return nil, fmt.Errorf("local transportations batch get (%s): %w", spreadsheetID, err)
		}

		var f []*sheets.DataFilter

		for _, s := range spr.Sheets {
			f = append(f, &sheets.DataFilter{
				A1Range: fmt.Sprintf("'%s'!A1:W", s.Properties.Title),
			})
		}

		b := sheets.BatchGetValuesByDataFilterRequest{DataFilters: f}

		resp, err := service.Spreadsheets.Values.BatchGetByDataFilter(spr.SpreadsheetId, &b).Do()
		if err != nil {
			return nil, fmt.Errorf("local transportations batch get (%s): %w", spreadsheetID, err)
		}

		rows := make([]map[string]string, 0)

		for _, v := range resp.ValueRanges {
			rows = append(rows, ExtractValuesFromRows(v.ValueRange, 1)...)

		}

		return rows, nil
	}
}

func GetAndMapRows[T any](rowGetter RowGetter, mapper RecordMapper[T], filters ...RowFilter) ([]T, error) {

	rows, err := rowGetter()

	if err != nil {
		return nil, err
	}
	mappedRows := make([]T, 0, len(rows))

	skip := func(row map[string]string) bool {
		for _, f := range filters {
			if !f(row) {
				return true
			}
		}
		return false
	}

	for i, row := range rows {

		if skip(row) {
			continue
		}

		mapped, err := mapper(row)
		if err != nil {
			fmt.Printf("error mapping in row %d: %v\n", i, err)
			continue
		}
		mappedRows = append(mappedRows, mapped)

	}

	return mappedRows, nil
}
