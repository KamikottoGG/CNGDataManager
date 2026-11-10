package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

type HotlineCallRecord struct {
	Timestamp           DateOnly       `json:"timestamp"`             // Позначка часу
	Responsible         string         `json:"responsible"`           // Відповідальний
	CallDate            DateOnly       `json:"call_date"`             // Дата дзвінка
	Region              NullString     `json:"region"`                // Область
	BeneficiaryFullName NullString     `json:"beneficiary_full_name"` // ПІБ людини, якій буде надана послуга
	Gender              BeneficiarySex `json:"gender"`                // Стать
	BirthDate           DateOnly       `json:"birth_date"`            // Дата народження
	Age                 int            `json:"age"`                   // Вік
	Phone               NullString     `json:"phone"`                 // Номер телефону
	Service             string         `json:"service"`               // Послуга
	Community           NullString     `json:"community"`             // Громада
	Settlement          NullString     `json:"settlement"`            // Населений пункт
	EvacuationAddress   NullString     `json:"evacuation_address"`    // Адреса евакуації
	Location            NullString     `json:"location"`              // Локація
	EvacuatedCount      int            `json:"evacuated_count"`       // Кількість евакуйованих
	ChildrenCount       int            `json:"children_count"`        // Кількість дітей
	MobileCount         int            `json:"mobile_count"`          // Кількість мобільних
	LowMobileCount      int            `json:"low_mobile_count"`      // Кількість маломобільних
	BedriddenCount      int            `json:"bedridden_count"`       // Кількість лежачих
	AnimalsCount        int            `json:"animals_count"`         // Кількість тварин
	LuggageCount        int            `json:"luggage_count"`         // Кількість речей
	NPEvacuation        NullString     `json:"np_evacuation"`         // НП евакуація (Нова Пошта или Населений Пункт)
	ApplicantInitiator  NullString     `json:"applicant_initiator"`   // Ініціатор(ка) заявки
	EvacuationDate      DateOnly       `json:"evacuation_date"`       // Дата евакуації
	OtherComment        NullString     `json:"other_comment"`         // Інший коментар
	Comment             NullString     `json:"comment"`               // Коментар
}

func (h *HotlineCallRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	parseInt := func(key string) int {
		val, _ := strconv.Atoi(get(key))
		return val
	}

	var err error

	h.Responsible = get("Відповідальний")
	h.BeneficiaryFullName = NullString(get("ПІБ людини, якій буде надана послуга"))
	h.Gender, err = ParseSex(get("Стать"))
	if err != nil {
		return fmt.Errorf("sex: %w", err)
	}

	h.Service = get("Послуга")

	if h.Timestamp, err = ParseSheetDate(get("Позначка часу")); err != nil {
		return fmt.Errorf("timestamp: %w", err)
	}

	if h.CallDate, err = ParseSheetDate(get("Дата дзвінка")); err != nil {
		return fmt.Errorf("call date: %w", err)
	}

	if h.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return fmt.Errorf("birth date: %w", err)
	}

	if h.EvacuationDate, err = ParseSheetDate(get("Дата евакуації")); err != nil {
		return fmt.Errorf("evacuation date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		h.Age = age
	} else if bt := time.Time(h.BirthDate); !bt.IsZero() {
		now := time.Now()
		h.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			h.Age--
		}
	}

	h.EvacuatedCount = parseInt("Кількість евакуйованих")
	h.ChildrenCount = parseInt("Кількість дітей")
	h.MobileCount = parseInt("Кількість мобільних")
	h.LowMobileCount = parseInt("Кількість маломобільних")
	h.BedriddenCount = parseInt("Кількість лежачих")
	h.AnimalsCount = parseInt("Кількість тварин")
	h.LuggageCount = parseInt("Кількість речей")

	h.Region = NullString(get("Область"))
	h.Phone = NullString(get("Номер телефону"))
	h.Community = NullString(get("Громада"))
	h.Settlement = NullString(get("Населений пункт"))
	h.EvacuationAddress = NullString(get("Адреса евакуації"))
	h.Location = NullString(get("Локація"))
	h.NPEvacuation = NullString(get("НП евакуація"))
	h.ApplicantInitiator = NullString(get("Ініціатор(ка) заявки"))
	h.OtherComment = NullString(get("Інший коментар"))
	h.Comment = NullString(get("Коментар"))

	return nil
}

func FilterHotlineCallValues(values []map[string]string) []map[string]string {
	newValues := make([]map[string]string, 0, len(values)/4)

	for _, value := range values {
		if len(strings.TrimSpace(value["Відповідальний"])) == 0 {
			continue
		}
		if len(strings.TrimSpace(value["Номер телефону"])) == 0 && len(strings.TrimSpace(value["Відповідальний"])) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return newValues
}

func CacheHotlineData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {

	log.Print("[+] ", "Hotline spreadsheet got with name: "+spreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))

	collectedData := make([]HotlineCallRecord, 0, len(spreadSheet.Sheets))

	rowsAffected := 0

	for i, sheet := range spreadSheet.Sheets {
		time.Sleep(1200 * time.Millisecond)
		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)

		rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

		if err != nil {
			return err
		}

		log.Printf("	[?] %d rows parsed. ", len(rows.Values))

		extractedValues := ExtractValuesFromRows(rows, 1)

		filteredValues := FilterHotlineCallValues(extractedValues)

		fvr := len(filteredValues)

		log.Printf("	[?] %d rows filtered. ", fvr)

		for j, v := range filteredValues {

			var mappedObject HotlineCallRecord

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

	log.Printf("[+] Hotline cache saved succesfully at path [%s]", jsonOutputFilePath)
	return nil
}
