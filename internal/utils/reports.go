package utils

import (
	"CNGDataManager/internal/google/sheets"
	"html"
	"strconv"
	"strings"
	"time"
)

func escape(s string) string {
	return html.EscapeString(s)
}

func bold(text string) string {
	return "<b>" + escape(text) + "</b>"
}

func writeGroupLine(sb *strings.Builder, label string, group GroupCount) {
	if group.Total == 0 {
		return
	}

	sb.WriteString("  - ")
	sb.WriteString(label)
	sb.WriteString(": <b>")
	sb.WriteString(strconv.Itoa(group.Total))
	sb.WriteString("</b>")

	if group.WithDisability > 0 {
		sb.WriteString(" <i>(інв. ")
		sb.WriteString(strconv.Itoa(group.WithDisability))
		sb.WriteString(")</i>")
	}
	sb.WriteByte('\n')
}

func writeSection(sb *strings.Builder, title string, data BeneficiariesDisaggregation, overwriteSecondLine string) {
	if title != "" {
		sb.WriteString(title)
		sb.WriteByte('\n')
	}

	if overwriteSecondLine != "" {
		sb.WriteString(overwriteSecondLine)
		sb.WriteString(": <b>")
	} else {
		sb.WriteString("Всього людей: <b>")
	}

	sb.WriteString(strconv.Itoa(data.Total()))
	sb.WriteString("</b>\n")

	writeGroupLine(sb, "Хлопці", data.Boys)
	writeGroupLine(sb, "Дівчата", data.Girls)
	writeGroupLine(sb, "Чоловіки", data.Men)
	writeGroupLine(sb, "Жінки", data.Women)
	writeGroupLine(sb, "Чол. похилого віку", data.ElderlyMen)
	writeGroupLine(sb, "Жін. похилого віку", data.ElderlyWomen)
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

	builder.WriteString("Звіт: ")
	builder.WriteString(startDate.Format("02.01.2006"))
	builder.WriteString(" - ")
	builder.WriteString(endDate.Format("02.01.2006"))
	builder.WriteString("\n\n")

	writeSection(&builder, "Евакуація:", evacuationDisaggregation, "")
	builder.WriteByte('\n')
	writeSection(&builder, "Прихисток:", shelterDisaggregation, "")
	builder.WriteByte('\n')
	writeSection(&builder, "Транзитний пункт:", transitDisaggregation, "")
	builder.WriteByte('\n')
	writeSection(&builder, "Перевезення:", localTransportationsDisaggregation, "")

	return strings.TrimRight(builder.String(), "\n")
}

func writeHeader(sb *strings.Builder, header string) {
	sb.WriteString(bold(header))
	sb.WriteByte('\n')
}

func writeReportHeader(sb *strings.Builder, startDate, endDate, cacheUpdateTime time.Time) {
	sb.WriteString("<b>Звіт CNG: ")

	format := "02.01"

	if startDate.Year() != endDate.Year() {
		format = "02.01.2006"
	}

	if startDate.Year() == 2000 {
		sb.WriteString("Увесь час")
	} else {
		sb.WriteString(startDate.Format(format))
		sb.WriteString(" - ")
		sb.WriteString(endDate.Format(format))
	}

	sb.WriteString("</b> (<i>оновл. ")
	sb.WriteString(cacheUpdateTime.Format("15:04"))
	sb.WriteString("</i>)\n\n")
}

func renderEvacuationSection(sb *strings.Builder, records []sheets.EvacuationRecord, tripRecords []sheets.VehicleTripRecord) {
	writeHeader(sb, "ЕВАКУАЦІЯ")

	evacuationTrips := 0
	for _, tripRecord := range tripRecords {
		if strings.TrimSpace(strings.ToLower(string(tripRecord.TripType))) != "евакуація" {
			continue
		}
		evacuationTrips++
	}

	sb.WriteString("Рейсів: <b>")
	sb.WriteString(strconv.Itoa(evacuationTrips))
	sb.WriteString("</b>\n")

	dis := GetEvacuationBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Евакуйовано")

	geoDis := GetEvacuationGeoDisaggregation(records)
	if len(geoDis) > 0 {
		sb.WriteString("Географія:\n")
		for r, hd := range geoDis {
			sb.WriteString("  <b>")
			sb.WriteString(escape(r))
			sb.WriteString(" обл.:</b> ")

			first := true
			for h, c := range hd {
				if !first {
					sb.WriteString(", ")
				}
				sb.WriteString(escape(h))
				sb.WriteString(" (<b>")
				sb.WriteString(strconv.Itoa(c))
				sb.WriteString("</b>)")
				first = false
			}
			sb.WriteByte('\n')
		}
	}

	sb.WriteByte('\n')
}

func renderMainShelterSection(sb *strings.Builder, records, allRecords []sheets.MainShelterRecord, startDate, endDate time.Time) {
	writeHeader(sb, "ОСНОВНИЙ ШЕЛТЕР")

	dis := GetMainShelterBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Прийнято")

	persResidents := make([]sheets.MainShelterRecord, 0)
	for _, record := range allRecords {
		if record.Status != "Проживає" {
			continue
		}
		persResidents = append(persResidents, record)
	}

	persDis := GetMainShelterBeneficiariesDisaggregation(persResidents)
	writeSection(sb, "", persDis, "Постійно проживає")

	depResidents := make([]sheets.MainShelterRecord, 0)

	for _, record := range allRecords {
		depDate := time.Time(record.DepartureDate)

		if depDate.IsZero() {
			continue
		}

		if isDateInRange(depDate, startDate, endDate) {
			depResidents = append(depResidents, record)
		}
	}
	depDis := GetMainShelterBeneficiariesDisaggregation(depResidents)
	writeSection(sb, "", depDis, "Виїхало")

	sb.WriteByte('\n')
}

func renderTransitHubSection(sb *strings.Builder, records []sheets.TransitHubRecord, allRecords []sheets.TransitHubRecord,
	startDate, endDate time.Time) {
	writeHeader(sb, "ТРАНЗИТНИЙ ХАБ")

	dis := GetTransitHubBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Прийнято")

	geoDis := GetTransitHubGeoDisaggregation(records)
	if len(geoDis) > 0 {
		sb.WriteString("Звідки:\n")
		for r, hd := range geoDis {
			sb.WriteString("  <b>")
			sb.WriteString(escape(r))
			sb.WriteString(" обл.:</b> ")

			first := true
			for h, c := range hd {
				if !first {
					sb.WriteString(", ")
				}
				sb.WriteString(escape(h))
				sb.WriteString(" (<b>")
				sb.WriteString(strconv.Itoa(c))
				sb.WriteString("</b>)")
				first = false
			}
			sb.WriteByte('\n')
		}
	}

	renderTransitHubRemainsSection(sb, allRecords)

	depRecords := make([]sheets.TransitHubRecord, 0)

	for _, record := range allRecords {

		depDate := time.Time(record.DepartureDate)
		if depDate.IsZero() {
			continue
		}

		if isDateInRange(depDate, startDate, endDate) {
			depRecords = append(depRecords, record)
		}
	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(depRecords), "Виїхало")

	sb.WriteByte('\n')
}

func renderTransitHubRemainsSection(sb *strings.Builder, records []sheets.TransitHubRecord) {
	fRecords := make([]sheets.TransitHubRecord, 0)
	for _, record := range records {
		if record.ResidenceStatus == "Заселений(а)" {
			fRecords = append(fRecords, record)
		}
	}

	writeSection(sb, "", GetTransitHubBeneficiariesDisaggregation(fRecords), "Залишилося")
}

func renderLocalTransportationsSection(sb *strings.Builder, records []sheets.LocalTransportationRecord) {
	writeHeader(sb, "ЛОКАЛЬНІ ПЕРЕВЕЗЕННЯ")
	dis := GetLocalTransportationsBeneficiariesDisaggregation(records)
	writeSection(sb, "", dis, "Перевезено")
	sb.WriteByte('\n')
}

func renderHotlineSection(sb *strings.Builder, records []sheets.HotlineCallRecord) {
	writeHeader(sb, "ГАРЯЧА ЛІНІЯ")

	if len(records) == 0 {
		sb.WriteString("Дзвінків не зафіксовано.\n\n")
		return
	}

	sb.WriteString("Всього дзвінків: <b>")
	sb.WriteString(strconv.Itoa(len(records)))
	sb.WriteString("</b>\n")

	callTypes := make(map[string]int)
	for _, record := range records {
		callTypes[record.Service]++
	}

	for s, c := range callTypes {
		sb.WriteString("  - ")
		sb.WriteString(escape(s))
		sb.WriteString(": <b>")
		sb.WriteString(strconv.Itoa(c))
		sb.WriteString("</b>\n")
	}

	sb.WriteByte('\n')
}

func isDateInRange(targetDate time.Time, start, end time.Time) bool {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, targetDate.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, targetDate.Location())
	return !targetDate.Before(startDay) && !targetDate.After(endDay)
}

func BuildFullReportText(
	reportData *ReportData,
) string {
	var sb strings.Builder
	sb.Grow(2048)

	writeReportHeader(&sb, reportData.StartDate, reportData.EndDate, reportData.CacheUpdateTime)
	renderEvacuationSection(&sb, reportData.Evacuation, reportData.VehicleTrip)
	renderMainShelterSection(&sb, reportData.Shelter, reportData.AllShelter, reportData.StartDate, reportData.EndDate)
	renderTransitHubSection(&sb, reportData.TransitHub, reportData.AllTransitHub, reportData.StartDate, reportData.EndDate)
	renderHotlineSection(&sb, reportData.Hotline)
	renderLocalTransportationsSection(&sb, reportData.LocalTransportations)

	return strings.TrimRight(sb.String(), "\n")
}
