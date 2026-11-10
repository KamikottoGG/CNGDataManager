package sheets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

type LocalTransportationRecord struct {
	Date            DateOnly            `json:"date"`             // Дата [ДД/ММ/РРРР]
	InformedConsent bool                `json:"informed_consent"` // Надаю інформо-вану згоду (так/ні)
	LastName        string              `json:"last_name"`        // Прізвище
	FirstName       string              `json:"first_name"`       // Ім'я
	MiddleName      NullString          `json:"middle_name"`      // По батькові
	BirthDate       DateOnly            `json:"birth_date"`       // Дата народження [ДД/ММ/РРРР]
	Age             int                 `json:"age"`              // Вік
	Sex             BeneficiarySex      `json:"sex"`              // Стать
	Disability      bool                `json:"disability"`       // Інвалідність (так/ні)
	Passport        BeneficiaryPassport `json:"passport"`         // Паспорт
	TaxID           NullString          `json:"tax_id"`           // ІПН
	FromRegion      NullString          `json:"from_region"`      // Звідки область
	FromCommunity   NullString          `json:"from_community"`   // Звідки Громада
	SettlementType  NullString          `json:"settlement_type"`  // Тип населеного пункту
	FromSettlement  NullString          `json:"from_settlement"`  // Звідки Наcпункт
	TransferredTo   NullString          `json:"transferred_to"`   // Кому передали
	TransportedTo   NullString          `json:"transported_to"`   // Куди перевезли
	Phone           NullString          `json:"phone"`            // Номер телефону
	Comment         NullString          `json:"comment"`          // Коментар
}

func (l *LocalTransportationRecord) MapObject(object map[string]string) error {
	get := func(key string) string { return strings.TrimSpace(object[key]) }

	var err error

	l.LastName = get("Прізвище")
	l.FirstName = get("Ім'я")

	if l.Sex, err = ParseSex(get("Стать")); err != nil {
		return fmt.Errorf("sex: %w", err)
	}

	l.InformedConsent = ParseYesNo(get("Надаю інформо-вану згоду (так/ні)"))
	l.Disability = ParseYesNo(get("Інвалідність (так/ні)"))

	if l.Date, err = ParseSheetDate(get("Дата [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("transportation date: %w", err)
	}

	if l.BirthDate, err = ParseSheetDate(get("Дата народження [ДД/ММ/РРРР]")); err != nil {
		return fmt.Errorf("birth date: %w", err)
	}

	if age, convErr := strconv.Atoi(get("Вік")); convErr == nil {
		l.Age = age
	} else if bt := time.Time(l.BirthDate); !bt.IsZero() {
		now := time.Now()
		l.Age = now.Year() - bt.Year()
		if now.YearDay() < bt.YearDay() {
			l.Age--
		}
	}

	l.Passport.Number = NullString(get("Паспорт"))

	l.MiddleName = NullString(get("По батькові"))
	l.TaxID = NullString(get("ІПН"))
	l.FromRegion = NullString(get("Звідки область"))
	l.FromCommunity = NullString(get("Звідки Громада"))
	l.SettlementType = NullString(get("Тип населеного пункту"))
	l.FromSettlement = NullString(get("Звідки Наcпункт"))
	l.TransferredTo = NullString(get("Кому передали"))
	l.TransportedTo = NullString(get("Куди перевезли"))
	l.Phone = NullString(get("Номер телефону"))
	l.Comment = NullString(get("Коментар"))

	return nil
}

func FilterLocalTransportationValues(values []map[string]string) []map[string]string {
	newValues := make([]map[string]string, 0, len(values)/4)

	for _, value := range values {
		if len(strings.TrimSpace(value["Прізвище"])) == 0 {
			continue
		}
		if len(strings.TrimSpace(value["Ім'я"])) == 0 {
			continue
		}

		newValues = append(newValues, value)
	}

	return newValues
}

func CacheLocalTransportationData(service *sheets.Service, spreadSheet *sheets.Spreadsheet, jsonOutputFilePath string) error {

	log.Print("[+] ", "Local transportation spreadsheet got with name: "+spreadSheet.Properties.Title)

	log.Printf("[?] Found %d sheets in spreadsheet", len(spreadSheet.Sheets))

	collectedData := make([]LocalTransportationRecord, 0, len(spreadSheet.Sheets))

	rowsAffected := 0

	for i, sheet := range spreadSheet.Sheets {
		log.Printf("[%d/%d] Starts work with sheet. Name: %s ", i+1, len(spreadSheet.Sheets), sheet.Properties.Title)

		rows, err := service.Spreadsheets.Values.Get(spreadSheet.SpreadsheetId, sheet.Properties.Title+"!A1:W").Do()

		if err != nil {
			return err
		}

		log.Printf("	[?] %d rows parsed. ", len(rows.Values))

		extractedValues := ExtractValuesFromRows(rows, 1)

		filteredValues := FilterLocalTransportationValues(extractedValues)

		fvr := len(filteredValues)

		log.Printf("	[?] %d rows filtered. ", fvr)

		for j, v := range filteredValues {

			var mappedObject LocalTransportationRecord

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

	log.Printf("[+] Local transportation cache saved succesfully at path [%s]", jsonOutputFilePath)
	return nil
}
