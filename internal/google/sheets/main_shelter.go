package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

type MainShelterRecord struct {
	AddedBy              string     `json:"added_by"`              // Внес дані
	ArrivalDate          DateOnly   `json:"arrival_date"`          // Дата прибуття (формат YYYY-MM-DD или DD.MM.YYYY)
	VPOStatus            string     `json:"vpo_status"`            // Статус ВПО
	IsTransit            bool       `json:"is_transit"`            // Транзит (true/false)
	LastName             string     `json:"last_name"`             // Прізвище
	FirstName            string     `json:"first_name"`            // Ім'я
	MiddleName           NullString `json:"middle_name"`           // По батькові (может отсутствовать -> null)
	BirthDate            DateOnly   `json:"birth_date"`            // Дата народження
	Age                  int        `json:"age"`                   // Вік
	Gender               string     `json:"gender"`                // Стать
	Category             string     `json:"category"`              // Категорія
	Disability           NullString `json:"disability"`            // Інвалідність (указатель, если null/нет)
	EvacuatedFromTown    NullString `json:"evacuated_from_town"`   // Населений пункт звідки евакуювали
	EvacuatedFromRegion  NullString `json:"evacuated_from_region"` // Область звідки евакуювали
	FurtherAccommodation NullString `json:"further_accommodation"` // Подальше розміщення
	TaxID                NullString `json:"tax_id"`                // ІПН (может быть пустым или содержать буквы для отказников)
	Phone                NullString `json:"phone"`                 // Телефон
	Passport             BeneficiaryPassport
	FamilyComposition    NullString `json:"family_composition"`   // Склад сім'ї
	RegistrationAddress  NullString `json:"registration_address"` // Адреса реєстрації
	Needs                NullString `json:"needs"`                // Потреби
	Status               string     `json:"status"`               // Статус (напр. "Проживает", "Выехал")
	DepartureDate        NullString `json:"departure_date"`       // Дата виїзду (если еще не выехал -> null)
}

func (m *MainShelterRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	var err error

	m.AddedBy = get("Внес дані")
	m.VPOStatus = get("Статус ВПО")
	m.LastName = get("Прізвище")
	m.FirstName = get("Ім'я")
	m.Gender = get("Стать")
	m.Category = get("Категорія")
	m.Status = get("Статус")

	m.IsTransit = ParseYesNo(get("Транзит"))

	if m.ArrivalDate, err = ParseSheetDate(get("Дата прибуття")); err != nil {
		return fmt.Errorf("arrival date: %w", err)
	}

	if m.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return fmt.Errorf("birth date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		m.Age = age
	} else if bt := time.Time(m.BirthDate); !bt.IsZero() {
		now := time.Now()
		m.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			m.Age--
		}
	}

	m.MiddleName = NullString(get("По батькові"))
	m.Disability = NullString(get("Інвалідність"))
	m.EvacuatedFromTown = NullString(get("Населений пункт звідки евакуювали"))
	m.EvacuatedFromRegion = NullString(get("Область звідки евакуювали"))
	m.FurtherAccommodation = NullString(get("Подальше розміщення"))
	m.TaxID = NullString(get("ІПН"))
	m.Phone = NullString(get("Телефон"))
	m.FamilyComposition = NullString(get("Склад сім'ї"))
	m.RegistrationAddress = NullString(get("Адреса реєстрації"))
	m.Needs = NullString(get("Потреби"))
	m.DepartureDate = NullString(get("Дата виїзду"))

	m.Passport.Number = NullString(get("Серія та номер паспорту"))
	m.Passport.Issuer = NullString(get("Ким виданий"))
	if m.Passport.IssueDate, err = ParseSheetDate(get("Коли виданий")); err != nil {
		return fmt.Errorf("passport issued date: %w", err)
	}

	return nil
}

func FilterMainShelterValues(values []map[string]string) []map[string]string {
	var newValues = make([]map[string]string, 0, len(values)/4)

	for _, value := range values {
		if len(value["Внес дані"]) == 0 {
			continue
		}
		if len(value["Ім'я"]) == 0 {
			continue
		}
		if len(value["Статус"]) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return newValues
}

func CacheMainShelterData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {

	log.Print("[+] ", "Main shelter spreadsheet got with name: "+spreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))

	rowsAffected := 0

	sheet, err := func() (*sheets.Sheet, error) {
		for _, sheet := range spreadSheet.Sheets {
			if sheet.Properties.Title == "Form_Responses" {
				return sheet, nil
			}
		}

		return nil, fmt.Errorf("could not find Form_Responses sheet")
	}()

	if err != nil {
		return err
	}

	log.Printf("[?] Starts work with sheet. Name: %s ", sheet.Properties.Title)

	rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

	if err != nil {
		return err
	}

	log.Printf("	[?] %d rows parsed. ", len(rows.Values))

	extractedValues := ExtractValuesFromRows(rows, 2)

	filteredValues := FilterMainShelterValues(extractedValues)

	fvr := len(filteredValues)

	log.Printf("	[?] %d rows filtered. ", fvr)

	mappedObjects := make([]MainShelterRecord, 0, fvr)

	for j, v := range filteredValues {

		var mappedObject MainShelterRecord

		if err := mappedObject.MapObject(v); err != nil {
			log.Printf("[!] row %d: %v", j, err)
			continue
		}

		mappedObjects = append(mappedObjects, mappedObject)

		rowsAffected++
	}

	log.Printf("[?] Work with sheet end. ")

	log.Printf("[+] %d rows mapped (not persisted yet)", rowsAffected)

	err = WriteToJsonFile(jsonOutputFilePath, mappedObjects)

	if err != nil {
		return err
	}

	log.Printf("[+] Main shelter cache saved succesfully at path [%s]", jsonOutputFilePath)
	return nil
}
