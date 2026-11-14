package sheets

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type EvacuationRecord struct {
	Date DateOnly `json:"date"`

	Name       string     `json:"name"`
	Surname    string     `json:"surname"`
	Patronymic NullString `json:"patronymic"`

	BirthDate      DateOnly            `json:"birthDate"`
	Age            int                 `json:"age"`
	Sex            BeneficiarySex      `json:"sex"`
	Disability     bool                `json:"disability"`
	Passport       BeneficiaryPassport `json:"passport"`
	TaxId          NullString          `json:"taxId"`
	Region         NullString          `json:"region"`
	District       NullString          `json:"district"`
	Community      NullString          `json:"community"`
	SettlementType NullString          `json:"settlementType"`
	Settlement     NullString          `json:"settlement"`

	EvacuatedTo NullString `json:"evacuatedTo"`
	PhoneNumber NullString `json:"phoneNumber"`
	Comment     NullString `json:"comment"`
}

//func (e *EvacuationRecord) IsEmpty() bool {
//	if time.Time(e.Date).IsZero() {
//		return false
//	}
//
//	if len(e.Surname) == 0 {
//		return false
//	}
//	if len(e.Name) == 0 {
//		return false
//	}
//
//	return true
//} func (e EvacuationRecord) IsEmpty(r map[string]string) bool {

func MapEvacuationRecord(object map[string]string) (EvacuationRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	var e EvacuationRecord
	var err error

	if e.Date, err = ParseSheetDate(get("Дата [ДД/ММ/РРРР]")); err != nil {
		return e, fmt.Errorf("evacuation date: %w", err)
	}

	e.Surname = get("Прізвище")
	e.Name = get("Ім'я")
	e.Patronymic = NullString(get("По батькові"))

	if e.BirthDate, err = ParseSheetDate(get("Дата народження [ДД/ММ/РРРР]")); err != nil {
		return e, fmt.Errorf("birth date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		e.Age = age
	} else if bt := time.Time(e.BirthDate); !bt.IsZero() {
		now := time.Now()
		e.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			e.Age--
		}
	}

	if e.Sex, err = ParseSex(get("Стать")); err != nil {
		return e, fmt.Errorf("sex: %w", err)
	}
	e.Disability = ParseYesNo(get("Інвалідність (так/ні)"))

	e.Passport.Number = NullString(get("Паспорт"))
	e.Passport.Issuer = NullString(get("Орган що видав"))
	if e.Passport.IssueDate, err = ParseSheetDate(get("Дата видачі")); err != nil {
		return e, fmt.Errorf("passport issue date: %w", err)
	}
	e.Passport.IssueCountry = NullString(get("Країна видачі"))

	e.TaxId = NullString(get("ІПН"))
	e.Region = NullString(get("Звідки область"))
	e.District = NullString(get("Звідки район"))
	e.Community = NullString(get("Звідки Громада"))
	e.SettlementType = NullString(get("Тип населеного пункту"))
	e.Settlement = NullString(get("Звідки Наcпункт"))

	e.EvacuatedTo = NullString(get("Куди евакуювалися"))
	e.PhoneNumber = NullString(get("Номер телефону"))
	e.Comment = NullString(get("Коментар"))

	return e, nil
}

//type EvacuationRowsCollector struct {
//	spreadsheet *sheets.Spreadsheet
//	service     *sheets.Service
//}
//
//func (e *EvacuationRowsCollector) Collect() ([]map[string]string, error) {
//	v := make([]map[string]string, 0)
//
//	for _, sheet := range e.spreadsheet.Sheets {
//		rows, err := e.service.Spreadsheets.Values.Get(e.spreadsheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()
//		if err != nil {
//			return nil, fmt.Errorf("get evacuation rows: %w", err)
//		}
//		extractedValues := ExtractValuesFromRows(rows, 1)
//
//	}
//	return v, nil
//}

//func NewEvacuationRecordsCollector(service *sheets.Service, spreadsheet *sheets.Spreadsheet) RecordsCollector[EvacuationRecord] {
//	return func(m RowMapper[EvacuationRecord], f ...RowFilter) ([]EvacuationRecord, error) {
//		records := make([]EvacuationRecord, 0)
//
//		for _, sheet := range spreadsheet.Sheets {
//			rows, err := service.Spreadsheets.Values.Get(spreadsheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()
//
//			if err != nil {
//				return nil, fmt.Errorf("get rows: %w", err)
//			}
//
//		}
//
//	}
//}

func FilterEvacuationValues(values []map[string]string) []map[string]string {
	var newValues = make([]map[string]string, 0, len(values)/4)

	for _, value := range values {
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

	return newValues

}

//func CacheEvacuationData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {
//
//	log.Print("[+] ", "Evacuation spreadsheet got with name: "+spreadSheet.Properties.Title)
//
//	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))
//
//	collectedData := make([]EvacuationRecord, 0, len(spreadSheet.Sheets))
//
//	rowsAffected := 0
//
//	for i, sheet := range spreadSheet.Sheets {
//		time.Sleep(1200 * time.Millisecond)
//		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)
//
//		rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()
//
//		if err != nil {
//			return err
//		}
//
//		log.Printf("	[?] %d rows parsed. ", len(rows.Values))
//
//		extractedValues := ExtractValuesFromRows(rows, 1)
//
//		filteredValues := FilterEvacuationValues(extractedValues)
//
//		fvr := len(filteredValues)
//
//		log.Printf("	[?] %d rows filtered. ", fvr)
//
//		for j, v := range filteredValues {
//
//			var mappedObject EvacuationRecord
//
//			if err := mappedObject.MapObject(v); err != nil {
//				log.Printf("[!] row %d: %v", j, err)
//				continue
//			}
//
//			collectedData = append(collectedData, mappedObject)
//
//			rowsAffected++
//		}
//
//		log.Printf("[%d/%d] Work with sheet end. ", i+1, len(spreadSheet.Sheets))
//
//	}
//
//	log.Printf("	[+] %d month iterated. ", len(collectedData))
//	log.Printf("[+] %d rows mapped (not persisted yet)", rowsAffected)
//
//	err := WriteToJsonFile(jsonOutputFilePath, collectedData)
//
//	if err != nil {
//		return err
//	}
//
//	log.Printf("[+] Evacuation cache saved succesfully at path [%s]", jsonOutputFilePath)
//	return nil
//}
