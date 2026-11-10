package utils

import (
	"CNGDataManager/internal/google/sheets"
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
	sb.WriteString(": <b>")
	sb.WriteString(strconv.Itoa(group.Total))
	sb.WriteString("</b>")

	if group.WithDisability > 0 {
		sb.WriteString(" <i>(з них з інвалідністю: ")
		sb.WriteString(strconv.Itoa(group.WithDisability))
		sb.WriteString(")</i>")
	}
	sb.WriteByte('\n')
}

func writeSection(sb *strings.Builder, title string, data BeneficiariesDisaggregation, overwriteSecondLine string) {
	if title != "" {
		sb.WriteString(title)
		sb.WriteString("\n")
	}

	if overwriteSecondLine != "" {
		sb.WriteString(overwriteSecondLine)
		sb.WriteString(": <b>")
	} else {
		sb.WriteString("Всього людей: <b>")
	}

	sb.WriteString(strconv.Itoa(data.Total()))
	sb.WriteString("</b>, з них:\n")

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
	builder.WriteString(" — ")
	builder.WriteString(endDate.Format("02.01.2006"))
	builder.WriteString("\n\n")

	writeSection(&builder, "Евакуація:", evacuationDisaggregation, "")
	builder.WriteString("\n")
	writeSection(&builder, "Прихисток:", shelterDisaggregation, "")
	builder.WriteString("\n")
	writeSection(&builder, "Транзитний пункт:", transitDisaggregation, "")
	builder.WriteString("\n")
	writeSection(&builder, "Перевезення:", localTransportationsDisaggregation, "")

	return strings.TrimRight(builder.String(), "\n")
}

func writeHeader(sb *strings.Builder, header string) {
	sb.WriteString("\n— — — — — — — — — — — — —\n")
	sb.WriteString(bold(header))
	sb.WriteString("\n\n")
}

func writeReportHeader(sb *strings.Builder, startDate, endDate, cacheUpdateTime time.Time) {
	sb.WriteString("<i>Кеш оновлено: ")
	sb.WriteString(cacheUpdateTime.Format("15:04:05"))
	sb.WriteString("</i>\n\n")

	sb.WriteString("<b>Звіт CNG за ")
	sb.WriteString(startDate.Format("02.01"))
	sb.WriteString(" — ")
	sb.WriteString(endDate.Format("02.01"))
	sb.WriteString("</b>\n")
	sb.WriteString("<i>Нижче розбивка по кожному проекту за цей проміжок часу</i>\n")
}

func bold(text string) string {
	return "<b>" + text + "</b>"
}

func renderEvacuationSection(sb *strings.Builder, records []sheets.EvacuationRecord, tripRecords []sheets.VehicleTripRecord) {
	writeHeader(sb, "ЕВАКУАЦІЯ")

	dis := GetEvacuationBeneficiariesDisaggregation(records)

	evacuationTrips := 0
	for _, tripRecord := range tripRecords {
		if strings.ToLower(string(tripRecord.TripType)) != "евакуація" {
			continue
		}
		evacuationTrips++
	}

	sb.WriteString("Кількість евакуаційних рейсів: <b>")
	sb.WriteString(strconv.Itoa(evacuationTrips))
	sb.WriteString("</b>\n")

	writeSection(sb, "", dis, "Всього евакуйовано")

	geoDis := GetEvacuationGeoDisaggregation(records)
	if len(geoDis) > 0 {
		sb.WriteString("\n<b>Географічний розподіл (область / громада):</b>\n")

		for r, hd := range geoDis {
			sb.WriteString(r)
			sb.WriteString(" обл.:\n")

			for h, c := range hd {
				sb.WriteString("  • ")
				sb.WriteString(h)
				sb.WriteString(": <b>")
				sb.WriteString(strconv.Itoa(c))
				sb.WriteString("</b>\n")
			}
		}
	}

	sb.WriteByte('\n')
}

func renderMainShelterSection(sb *strings.Builder, records, allRecords []sheets.MainShelterRecord) {
	writeHeader(sb, "ОСНОВНИЙ ШЕЛТЕР")

	dis := GetMainShelterBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Прийнято людей")
	sb.WriteByte('\n')

	persResidents := make([]sheets.MainShelterRecord, 0)
	for _, record := range allRecords {
		if record.Status != "Проживає" {
			continue
		}
		persResidents = append(persResidents, record)
	}

	persDis := GetMainShelterBeneficiariesDisaggregation(persResidents)
	writeSection(sb, "", persDis, "Постійно проживає людей")

	sb.WriteByte('\n')
}

func renderTransitHubSection(sb *strings.Builder, records []sheets.TransitHubRecord, allRecords []sheets.TransitHubRecord) {
	writeHeader(sb, "ТРАНЗИТНИЙ ХАБ")

	dis := GetTransitHubBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Людей прийнято")

	geoDis := GetTransitHubGeoDisaggregation(records)
	if len(geoDis) > 0 {
		sb.WriteString("\n<b>Прийняті з (область / громада):</b>\n")

		for r, hd := range geoDis {
			sb.WriteString(r)
			sb.WriteString(" обл.:\n")

			for h, c := range hd {
				sb.WriteString("  • ")
				sb.WriteString(h)
				sb.WriteString(": <b>")
				sb.WriteString(strconv.Itoa(c))
				sb.WriteString("</b>\n")
			}
		}
	}

	sb.WriteByte('\n')
	renderTransitHubRemainsSection(sb, allRecords)
	sb.WriteByte('\n')

	depRecords := make([]sheets.TransitHubRecord, 0)
	for _, record := range records {
		if record.ResidenceStatus == "Заселений(а)" {
			continue
		}
		depRecords = append(depRecords, record)
	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(depRecords), "Людей виїхало")

	sb.WriteByte('\n')
}

func renderTransitHubRemainsSection(sb *strings.Builder, records []sheets.TransitHubRecord) {
	fRecords := make([]sheets.TransitHubRecord, 0)
	for _, record := range records {
		if record.ResidenceStatus == "Заселений(а)" {
			fRecords = append(fRecords, record)
		}
	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(fRecords), "Людей залишилося")
}

func renderLocalTransportationsSection(sb *strings.Builder, records []sheets.LocalTransportationRecord) {
	writeHeader(sb, "ЛОКАЛЬНІ ПЕРЕВЕЗЕННЯ")
	dis := GetLocalTransportationsBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Людей перевезено")
	sb.WriteByte('\n')
}

func renderHotlineSection(sb *strings.Builder, records []sheets.HotlineCallRecord) {
	writeHeader(sb, "ГАРЯЧА ЛІНІЯ")

	if len(records) == 0 {
		sb.WriteString("Жодного дзвінка за період не зафіксовано.\n\n")
		return
	}

	sb.WriteString("Дзвінків за період: <b>")
	sb.WriteString(strconv.Itoa(len(records)))
	sb.WriteString("</b>\n")
	sb.WriteString("З них:\n")

	callTypes := make(map[string]int)
	for _, record := range records {
		callTypes[record.Service]++
	}

	for s, c := range callTypes {
		sb.WriteString("  • ")
		sb.WriteString(s)
		sb.WriteString(": <b>")
		sb.WriteString(strconv.Itoa(c))
		sb.WriteString("</b>\n")
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
	vehicleTripRecords []sheets.VehicleTripRecord,
	cacheUpdateTime time.Time,
) string {
	var sb strings.Builder
	sb.Grow(2048)

	writeReportHeader(&sb, startDate, endDate, cacheUpdateTime)
	renderEvacuationSection(&sb, evacuationRecords, vehicleTripRecords)
	renderMainShelterSection(&sb, shelterRecords, allShelterRecords)
	renderTransitHubSection(&sb, transitRecords, allTransitRecords)
	renderHotlineSection(&sb, hotlineRecords)
	renderLocalTransportationsSection(&sb, transportationsDisaggregation)

	return strings.TrimRight(sb.String(), "\n")
}
