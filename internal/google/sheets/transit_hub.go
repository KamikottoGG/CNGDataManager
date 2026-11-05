package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

// TransitHubRecord описывает запись в транзитном хабе
type TransitHubRecord struct {
	ResidenceStatus        string              `json:"residence_status"`        // Статус проживання
	ArrivalDate            DateOnly            `json:"arrival_date"`            // Дата прибуття
	DepartureDate          DateOnly            `json:"departure_date"`          // Дата вибуття
	FromRegion             NullString          `json:"from_region"`             // Звідки область
	FromDistrict           NullString          `json:"from_district"`           // Звідки район
	FromCommunity          NullString          `json:"from_community"`          // Звідки громада
	FromSettlement         NullString          `json:"from_settlement"`         // Звідки нас пункт
	LastName               string              `json:"last_name"`               // Прізвище
	FirstName              string              `json:"first_name"`              // Ім'я
	MiddleName             NullString          `json:"middle_name"`             // По батькові
	Gender                 string              `json:"gender"`                  // Стать
	BirthDate              DateOnly            `json:"birth_date"`              // Дата народження
	Age                    int                 `json:"age"`                     // Вік
	TaxID                  NullString          `json:"tax_id"`                  // ІПН
	Passport               BeneficiaryPassport `json:"passport"`                // Паспортні данные
	RegistrationRegion     NullString          `json:"registration_region"`     // Прописка: Область
	RegistrationCommunity  NullString          `json:"registration_community"`  // Прописка: Громада
	RegistrationSettlement NullString          `json:"registration_settlement"` // Прописка: Нас пункт
	RegistrationAddress    NullString          `json:"registration_address"`    // Прописка: Адреса
	Phone                  NullString          `json:"phone"`                   // Телефон
	Disability             NullString          `json:"disability"`              // Інвалідність (так/ні)
	Category               NullString          `json:"category"`                // Категорія/коротка інформація про осіб
	FurtherAccommodation   NullString          `json:"further_accommodation"`   // Подальше розміщення
	PlannedDepartureDate   DateOnly            `json:"planned_departure_date"`  // Планова дата виїзду
	Needs                  NullString          `json:"needs"`                   // Потреби
}

func (t *TransitHubRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	var err error

	t.ResidenceStatus = get("Статус проживання")
	t.LastName = get("Прізвище")
	t.FirstName = get("Ім'я")
	t.Gender = get("Стать")

	if t.ArrivalDate, err = ParseSheetDate(get("Дата прибуття")); err != nil {
		return fmt.Errorf("arrival date: %w", err)
	}

	if t.DepartureDate, err = ParseSheetDate(get("Дата вибуття")); err != nil {
		return fmt.Errorf("departure date: %w", err)
	}

	if t.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return fmt.Errorf("birth date: %w", err)
	}

	if t.PlannedDepartureDate, err = ParseSheetDate(get("Планова дата виїзду")); err != nil {
		return fmt.Errorf("planned departure date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		t.Age = age
	} else if bt := time.Time(t.BirthDate); !bt.IsZero() {
		now := time.Now()
		t.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			t.Age--
		}
	}

	t.Passport.Number = NullString(get("Паспорт/ID/Свідоцтво про народження"))
	t.Passport.Issuer = NullString(get("ким виданий"))
	if t.Passport.IssueDate, err = ParseSheetDate(get("дата видачі")); err != nil {
		return fmt.Errorf("passport issue date: %w", err)
	}

	t.FromRegion = NullString(get("Звідки область"))
	t.FromDistrict = NullString(get("Звідки район"))
	t.FromCommunity = NullString(get("Звідки громада"))
	t.FromSettlement = NullString(get("Звідки нас пункт"))

	t.MiddleName = NullString(get("По батькові"))
	t.TaxID = NullString(get("ІПН"))

	t.RegistrationRegion = NullString(get("Прописка: Область"))
	t.RegistrationCommunity = NullString(get("Прописка: Громада"))
	t.RegistrationSettlement = NullString(get("Прописка: Нас пункт"))
	t.RegistrationAddress = NullString(get("Прописка: Адреса"))

	t.Phone = NullString(get("Телефон"))
	t.Disability = NullString(get("Інвалідність (так/ні)"))
	t.Category = NullString(get("Категорія/коротка інформація про осіб"))
	t.FurtherAccommodation = NullString(get("Подальше розміщення"))
	t.Needs = NullString(get("Потреби"))

	return nil
}

func FilterTransitHubValues(values []map[string]string) []map[string]string {
	newValues := make([]map[string]string, 0, len(values)/4)

	for _, value := range values {
		if len(strings.TrimSpace(value["Прізвище"])) == 0 {
			continue
		}
		if len(strings.TrimSpace(value["Ім'я"])) == 0 {
			continue
		}

		if len(strings.TrimSpace(value["Статус проживання"])) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return newValues
}

func CacheTransitHubData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {

	log.Print("[+] ", "Transit hub spreadsheet got with name: "+spreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))

	collectedData := make([]TransitHubRecord, 0, len(spreadSheet.Sheets))

	rowsAffected := 0

	for i, sheet := range spreadSheet.Sheets {
		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)

		rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

		if err != nil {
			return err
		}

		log.Printf("	[?] %d rows parsed. ", len(rows.Values))

		extractedValues := ExtractValuesFromRows(rows, 1)

		filteredValues := FilterTransitHubValues(extractedValues)

		fvr := len(filteredValues)

		log.Printf("	[?] %d rows filtered. ", fvr)

		for j, v := range filteredValues {

			var mappedObject TransitHubRecord

			if err := mappedObject.MapObject(v); err != nil {
				log.Printf("[!] row %d: %v", j, err)
				continue
			}

			collectedData = append(collectedData, mappedObject)

			rowsAffected++
		}

		log.Printf("[%d/%d] Work with sheet end. ", i+1, len(spreadSheet.Sheets))

	}

	log.Printf("	[+] %d month iterated. ", len(collectedData))
	log.Printf("[+] %d rows mapped (not persisted yet)", rowsAffected)

	err := WriteToJsonFile(jsonOutputFilePath, collectedData)

	if err != nil {
		return err
	}

	log.Printf("[+] Transit hub cache saved succesfully at path [%s]", jsonOutputFilePath)
	return nil
}
