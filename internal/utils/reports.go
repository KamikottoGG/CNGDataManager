package utils

import (
	"CNGDataManager/internal/google/sheets"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func writeGroupLine(sb *strings.Builder, label string, group GroupCount) {
	if group.Total == 0 {
		return
	}

	sb.WriteString("  • ")
	sb.WriteString(label)
	sb.WriteString(": ")
	sb.WriteString(strconv.Itoa(group.Total))

	if group.WithDisability > 0 {
		sb.WriteString(" (з них з інвалідністю: ")
		sb.WriteString(strconv.Itoa(group.WithDisability))
		sb.WriteString(")")
	}
	sb.WriteByte('\n')
}

func writeSection(sb *strings.Builder, title string, data BeneficiariesDisaggregation, overwriteSecondLine string) {
	sb.WriteString(title)
	if overwriteSecondLine != "" {
		sb.WriteString("\n")
		sb.WriteString(overwriteSecondLine)
	} else {
		sb.WriteString("\nВсього людей: ")
	}

	sb.WriteString(strconv.Itoa(data.Total()))

	sb.WriteString(", з них")
	sb.WriteByte('\n')

	writeGroupLine(sb, "Хлопці", data.Boys)
	writeGroupLine(sb, "Дівчата", data.Girls)
	writeGroupLine(sb, "Чоловіки", data.Men)
	writeGroupLine(sb, "Жінки", data.Women)
	writeGroupLine(sb, "Чоловіки похилого віку", data.ElderlyMen)
	writeGroupLine(sb, "Жінки похилого віку", data.ElderlyWomen)

}

func BuildReportTextRaw(
	startDate time.Time,
	endDate time.Time,
	evacuationDisaggregation,
	shelterDisaggregation,
	transitDisaggregation,
	localTransportationsDisaggregation BeneficiariesDisaggregation,
) string {
	var builder strings.Builder
	builder.Grow(1024)

	builder.WriteString("Звіт за ")
	builder.WriteString(startDate.Format("02.01.2006"))
	builder.WriteString(" -> ")
	builder.WriteString(endDate.Format("02.01.2006"))
	builder.WriteString("\n\n")

	writeSection(&builder, "Евакуація:", evacuationDisaggregation, "")
	writeSection(&builder, "Прихисток:", shelterDisaggregation, "")
	writeSection(&builder, "Транзитний пункт:", transitDisaggregation, "")
	writeSection(&builder, "Перевезення:", localTransportationsDisaggregation, "")

	return strings.TrimRight(builder.String(), "\n")
}

func writeHeader(sb *strings.Builder, header string) {

	sb.WriteString(bold(header))

	sb.WriteString("\n")
}

func writeReportHeader(sb *strings.Builder, startDate, endDate, cacheUpdateTime time.Time) {

	sb.WriteString("<i>")

	sb.WriteString("Кеш оновлено: ")
	sb.WriteString(cacheUpdateTime.Format("15:04:05"))
	sb.WriteString("</i>\n\n")

	sb.WriteString("Звіт CNG за ")
	sb.WriteString(startDate.Format("02.01"))
	sb.WriteString(" - ")
	sb.WriteString(endDate.Format("02.01"))
	sb.WriteString("\n")
	sb.WriteString("Нижче розбивка по кожному проекту за цей проміжок часу \n\n")

}

func bold(text string) string {
	return "<b>" + text + "</b>"
}

func renderEvacuationSection(sb *strings.Builder, records []sheets.EvacuationRecord) {
	writeHeader(sb, "ЕВАКУАЦІЯ")

	dis := GetEvacuationBeneficiariesDisaggregation(records)

	sb.WriteString("Кількість евакуаційних рейсів: NaN")
	writeSection(sb, "", dis, "Всього евакуйовано ")

	geoDis := GetEvacuationGeoDisaggregation(records)

	sb.WriteString("Географічний розподіл(область/громада)\n")

	for r, hd := range geoDis {
		sb.WriteString(bold(r))
		sb.WriteByte('\n')

		for h, c := range hd {
			sb.WriteString("  • ")
			sb.WriteString(h)
			sb.WriteString(": ")
			sb.WriteString(strconv.Itoa(c))
			sb.WriteByte('\n')
		}

		sb.WriteByte('\n')

	}

	sb.WriteByte('\n')
}

func renderMainShelterSection(sb *strings.Builder, records, allRecords []sheets.MainShelterRecord) {
	sb.WriteString(bold("ОСНОВНИЙ ШЕЛТЕР"))
	dis := GetMainShelterBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Прийнято людей ")

	persResidents := make([]sheets.MainShelterRecord, 0)

	for _, record := range allRecords {

		if record.Status != "Проживає" {
			continue
		}

		// TODO: Доделать фильтр, уточнить у Ири за каким критерием фильтрація
		persResidents = append(persResidents, record)

	}

	persDis := GetMainShelterBeneficiariesDisaggregation(persResidents)

	writeSection(sb, "", persDis, "Постійно проживає людей ")

	sb.WriteByte('\n')

}

func renderTransitHubSection(sb *strings.Builder, records []sheets.TransitHubRecord, allRecords []sheets.TransitHubRecord) {

	dis := GetTransitHubBeneficiariesDisaggregation(records)

	writeSection(sb, "<b>ТРАНЗИТНИЙ ХАБ</b>", dis, "Людей прийнято ")

	geoDis := GetTransitHubGeoDisaggregation(records)

	sb.WriteString("Прийняті з(область/громада)\n")

	for r, hd := range geoDis {
		sb.WriteString(bold(r))
		sb.WriteByte('\n')

		for h, c := range hd {
			sb.WriteString("  • ")
			sb.WriteString(h)
			sb.WriteString(": ")
			sb.WriteString(strconv.Itoa(c))
			sb.WriteByte('\n')
		}

		sb.WriteByte('\n')

	}

	renderTransitHubRemainsSection(sb, allRecords)

	depRecords := make([]sheets.TransitHubRecord, 0)

	for _, record := range records {
		if record.ResidenceStatus == "Заселений(а)" {
			continue
		}

		depRecords = append(depRecords, record)
	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(depRecords), "Людей виїхало ")

	sb.WriteByte('\n')
}

func renderTransitHubRemainsSection(sb *strings.Builder, records []sheets.TransitHubRecord) {
	fRecords := make([]sheets.TransitHubRecord, 0)

	for _, record := range records {
		if record.ResidenceStatus == "Заселений(а)" {
			fRecords = append(fRecords, record)
		}

	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(fRecords), "Людей залишилося ")

}

func renderLocalTransportationsSection(sb *strings.Builder, records []sheets.LocalTransportationRecord) {
	dis := GetLocalTransportationsBeneficiariesDisaggregation(records)
	writeSection(sb, "<b>ЛОКАЛЬНІ ПЕРЕВЕЗЕННЯ</b>", dis, "Людей перевезено ")

}

func renderHotlineSection(sb *strings.Builder, records []sheets.HotlineCallRecord) {

	writeHeader(sb, "ГАРЯЧА ЛІНІЯ")

	if len(records) == 0 {
		sb.WriteString("Жодного дзвінка за період не зафіксовано.")
		return
	}

	sb.WriteString("Дзвінків за період: ")
	sb.WriteString(strconv.Itoa(len(records)))
	sb.WriteString("\nЗ них\n")

	callTypes := make(map[string]int, 0)
	fmt.Println(records)
	for _, record := range records {
		callTypes[record.Service]++
		fmt.Println(record.Service)
	}

	for s, c := range callTypes {
		sb.WriteString("  • ")
		sb.WriteString(s)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(c))
		sb.WriteByte('\n')
	}

	sb.WriteByte('\n')
}

func BuildFullReportText(
	startDate time.Time,

	endDate time.Time,
	hotlineRecords []sheets.HotlineCallRecord,
	evacuationRecords []sheets.EvacuationRecord,
	shelterRecords []sheets.MainShelterRecord,
	allShelterRecords []sheets.MainShelterRecord,
	transitRecords []sheets.TransitHubRecord,
	allTransitRecords []sheets.TransitHubRecord,
	transportationsDisaggregation []sheets.LocalTransportationRecord,
	cacheUpdateTime time.Time,
) string {

	sb := strings.Builder{}

	sb.Grow(1024)

	writeReportHeader(&sb, startDate, endDate, cacheUpdateTime)

	renderEvacuationSection(&sb, evacuationRecords)

	renderMainShelterSection(&sb, shelterRecords, allShelterRecords)

	renderTransitHubSection(&sb, transitRecords, allTransitRecords)

	renderHotlineSection(&sb, hotlineRecords)

	renderLocalTransportationsSection(&sb, transportationsDisaggregation)

	return sb.String()
}
