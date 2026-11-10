package utils

import (
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

func writeSection(sb *strings.Builder, title string, data BeneficiariesDisaggregation) {
	sb.WriteString(title)
	sb.WriteString("\nВсього: ")
	sb.WriteString(strconv.Itoa(data.Total()))
	sb.WriteByte('\n')

	writeGroupLine(sb, "Хлопці (<18)", data.Boys)
	writeGroupLine(sb, "Дівчата (<18)", data.Girls)
	writeGroupLine(sb, "Чоловіки (18-59)", data.Men)
	writeGroupLine(sb, "Жінки (18-59)", data.Women)
	writeGroupLine(sb, "Чоловіки похилого віку (60+)", data.ElderlyMen)
	writeGroupLine(sb, "Жінки похилого віку (60+)", data.ElderlyWomen)
	sb.WriteByte('\n')
}

func BuildReportTextRaw(
	startDate time.Time,
	endDate time.Time,
	evacuationDisaggregation,
	shelterDisaggregation,
	transitDisaggregation,
	transportationsDisaggregation BeneficiariesDisaggregation,
) string {
	var builder strings.Builder
	builder.Grow(1024)

	builder.WriteString("Звіт за ")
	builder.WriteString(startDate.Format("02.01.2006"))
	builder.WriteString(" -> ")
	builder.WriteString(endDate.Format("02.01.2006"))
	builder.WriteString("\n\n")

	writeSection(&builder, "Евакуація:", evacuationDisaggregation)
	writeSection(&builder, "Прихисток:", shelterDisaggregation)
	writeSection(&builder, "Транзитний пункт:", transitDisaggregation)
	writeSection(&builder, "Перевезення:", transportationsDisaggregation)

	return strings.TrimRight(builder.String(), "\n")
}
