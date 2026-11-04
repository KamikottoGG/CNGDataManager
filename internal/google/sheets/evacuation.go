package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
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

func (e *EvacuationRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	var err error

	if e.Date, err = ParseSheetDate(get("Дата [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("evacuation date: %w", err)
	}

	e.Surname = get("Прізвище")
	e.Name = get("Ім'я")
	e.Patronymic = NullString(get("По батькові"))

	if e.BirthDate, err = ParseSheetDate(get("Дата народження [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("birth date: %w", err)
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
		return fmt.Errorf("sex: %w", err)
	}
	e.Disability = ParseYesNo(get("Інвалідність (так/ні)"))

	e.Passport.Number = get("Паспорт")
	e.Passport.Issuer = get("Орган що видав")
	if e.Passport.IssueDate, err = ParseSheetDate(get("Дата видачі")); err != nil {
		return fmt.Errorf("passport issue date: %w", err)
	}
	e.Passport.IssueCountry = get("Країна видачі")

	e.TaxId = NullString(get("ІПН"))
	e.Region = NullString(get("Звідки область"))
	e.District = NullString(get("Звідки район"))
	e.Community = NullString(get("Звідки Громада"))
	e.SettlementType = NullString(get("Тип населеного пункту"))
	e.Settlement = NullString(get("Звідки Наcпункт"))

	e.EvacuatedTo = NullString(get("Куди евакуювалися"))
	e.PhoneNumber = NullString(get("Номер телефону"))
	e.Comment = NullString(get("Коментар"))

	return nil
}

func CacheEvacuationData(service *sheets.Service, evacuationSpreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {

	log.Print("[+] ", "Evacuation spreadsheet got with name: "+evacuationSpreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(evacuationSpreadSheet.Sheets))

	collectedData := make([]*[]EvacuationRecord, 0, len(evacuationSpreadSheet.Sheets))

	rowsAffected := 0

	for i, sheet := range evacuationSpreadSheet.Sheets {
		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(evacuationSpreadSheet.Sheets), sheet.Properties.Title)

		rows, err := service.Spreadsheets.Values.Get(evacuationSpreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

		if err != nil {
			return err
		}

		log.Printf("	[?] %d rows parsed. ", len(rows.Values))

		extractedValues := ExtractValuesFromRows(rows)

		filteredValues := FilterValues(extractedValues)

		fvr := len(filteredValues)

		log.Printf("	[?] %d rows filtered. ", fvr)

		mappedObjects := make([]EvacuationRecord, 0, fvr)

		for j, v := range filteredValues {

			var mappedObject EvacuationRecord

			if err := mappedObject.MapObject(v); err != nil {
				log.Printf("[!] row %d: %v", j, err)
				continue
			}

			mappedObjects = append(mappedObjects, mappedObject)

			rowsAffected++
		}

		collectedData = append(collectedData, &mappedObjects)

		log.Printf("[%d/%d] Work with sheet end. ", i+1, len(evacuationSpreadSheet.Sheets))

	}

	log.Printf("	[+] %d month iterated. ", len(collectedData))
	log.Printf("[+] %d rows mapped (not persisted yet)", rowsAffected)

	err := WriteToJsonFile(jsonOutputFilePath, collectedData)

	if err != nil {
		return err
	}

	log.Printf("[+] Evacuation cache saved succesfully at path [%s]", jsonOutputFilePath)
	return nil
}
