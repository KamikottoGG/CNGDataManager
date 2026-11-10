package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

type VehicleTripRecord struct {
	Date DateOnly `json:"date"`

	StartTime NullString `json:"startTime"` // СТАРТ
	EndTime   NullString `json:"endTime"`   // ФІНІШ

	Driver       string     `json:"driver"`       // Водій/відповідальний за виїзд
	Team         NullString `json:"team"`         // Команда виїзду
	VehiclePlate string     `json:"vehiclePlate"` // Марка машини – номери
	Route        string     `json:"route"`        // Маршрут
	TripType     NullString `json:"tripType"`     // Тип виїзду

	EvacueesCount    int `json:"evacueesCount"`    // Кількість евакуйованих осіб
	LowMobilityCount int `json:"lowMobilityCount"` // Кількість інвалідів І маломобільних
	ChildrenCount    int `json:"childrenCount"`    // Кількість дітей

	OdometerStart float64 `json:"odometerStart"` // Показання одометра на початку поїздки, км
	OdometerEnd   float64 `json:"odometerEnd"`   // Показання одометра на кінець поїздки, км
	FuelRefueled  float64 `json:"fuelRefueled"`  // Заправлено палива, л
	FuelRemaining float64 `json:"fuelRemaining"` // Залишок палива(л.)

	Comments NullString `json:"comments"`   // Коментарі
	ExtraCol NullString `json:"extraCol16"` // Стовпець 16 (технічна або запасна колонка)
}

func (v *VehicleTripRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	parseOptionalInt := func(val string) int {
		clean := strings.TrimSpace(val)
		if clean == "" {
			return 0
		}
		num, err := strconv.Atoi(clean)
		if err != nil {
			return 0
		}
		return num
	}

	parseOptionalFloat := func(val string) float64 {
		clean := strings.TrimSpace(strings.ReplaceAll(val, ",", "."))
		if clean == "" {
			return 0.0
		}
		f, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			return 0.0
		}
		return f
	}

	var err error

	//tsRaw := get("Позначка часу")
	//if tsRaw != "" {
	//
	//	if parsedTime, parseErr := time.Parse("02.01.2006 15:04:05", tsRaw); parseErr == nil {
	//		v.Timestamp = parsedTime
	//	} else if parsedTime, parseErr = time.Parse("2006-01-02 15:04:05", tsRaw); parseErr == nil {
	//		v.Timestamp = parsedTime
	//	}
	//}

	if v.Date, err = ParseSheetDate(get("Дата")); err != nil {
		return fmt.Errorf("trip date: %w", err)
	}

	v.StartTime = NullString(get("СТАРТ"))
	v.EndTime = NullString(get("ФІНІШ"))

	v.Driver = get("Водій/відповідальний за виїзд")
	v.Team = NullString(get("Команда виїзду"))
	v.VehiclePlate = get("Марка машини – номери")
	v.Route = get("Маршрут")
	v.TripType = NullString(get("Тип виїзду"))

	v.EvacueesCount = parseOptionalInt(get("Кількість евакуйованих осіб"))
	v.LowMobilityCount = parseOptionalInt(get("Кількість інвалідів І маломобільних"))
	v.ChildrenCount = parseOptionalInt(get("Кількість дітей"))

	v.OdometerStart = parseOptionalFloat(get("Показання oдометра на початку поїздки, км"))
	if v.OdometerStart == 0 {
		v.OdometerStart = parseOptionalFloat(get("Показання oдометра на\u00a0початкупоїздки, км"))
	}

	v.OdometerEnd = parseOptionalFloat(get("Показання oдометра на кінець поїздки, км"))
	if v.OdometerEnd == 0 {
		v.OdometerEnd = parseOptionalFloat(get("Показання oдометра на\u00a0кiнецьпоїздки, км"))
	}

	v.FuelRefueled = parseOptionalFloat(get("Заправлено палива, л"))
	v.FuelRemaining = parseOptionalFloat(get("Залишок палива(л.)"))

	v.Comments = NullString(get("Коментарі"))
	v.ExtraCol = NullString(get("Стовпець 16"))

	return nil
}

func FilterVehicleTripValues(values []map[string]string) []map[string]string {
	newValues := make([]map[string]string, 0, len(values))

	for _, value := range values {

		if len(strings.TrimSpace(value["Дата"])) == 0 && len(strings.TrimSpace(value["Позначка часу"])) == 0 {
			continue
		}
		if len(strings.TrimSpace(value["Водій/відповідальний за виїзд"])) == 0 {
			continue
		}
		if len(strings.TrimSpace(value["Марка машини – номери"])) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return newValues
}

func CacheVehicleTripData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {
	log.Printf("[+] Vehicle trip spreadsheet: %s", spreadSheet.Properties.Title)
	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))

	collectedData := make([]VehicleTripRecord, 0, 50)
	rowsAffected := 0

	for i, sheet := range spreadSheet.Sheets {

		time.Sleep(1200 * time.Millisecond)
		log.Printf("[%d/%d] Processing sheet: %s", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)
		if strings.Contains(sheet.Properties.Title, "Відповіді") {
			log.Printf("   [->] Skipped sheet %s", sheet.Properties.Title)
			continue
		}
		rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:R").Do()
		if err != nil {
			return fmt.Errorf("fetch rows for sheet %s: %w", sheet.Properties.Title, err)
		}

		log.Printf("   [?] %d rows found in sheet", len(rows.Values))

		extractedValues := ExtractValuesFromRows(rows, 1)
		filteredValues := FilterVehicleTripValues(extractedValues)

		log.Printf("   [?] %d valid rows after filter", len(filteredValues))

		for j, v := range filteredValues {
			var mappedObject VehicleTripRecord

			if err := mappedObject.MapObject(v); err != nil {
				log.Printf("[!] row %d mapping failed: %v", j, err)
				continue
			}

			collectedData = append(collectedData, mappedObject)
			rowsAffected++
		}

		log.Printf("[%d/%d] Finished sheet: %s", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)
	}

	log.Printf("[+] Total rows mapped: %d", rowsAffected)

	if err := WriteToJsonFile(jsonOutputFilePath, collectedData); err != nil {
		return fmt.Errorf("write vehicle trips json: %w", err)
	}

	log.Printf("[+] Vehicle trips cache saved successfully: [%s]", jsonOutputFilePath)
	return nil
}
