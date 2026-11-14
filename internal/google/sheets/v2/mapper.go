package sheets

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EvacuationMapper(object map[string]string) (EvacuationRecord, error) {
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

func TransitHubMapper(object map[string]string) (TransitHubRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	var t TransitHubRecord
	var err error

	t.ResidenceStatus = get("Статус проживання")
	t.LastName = get("Прізвище")
	t.FirstName = get("Ім'я")

	if t.Sex, err = ParseSex(get("Стать")); err != nil {
		return t, fmt.Errorf("sex: %w", err)
	}

	if t.ArrivalDate, err = ParseSheetDate(get("Дата прибуття")); err != nil {
		return t, fmt.Errorf("arrival date: %w", err)
	}

	if t.DepartureDate, err = ParseSheetDate(get("Дата вибуття")); err != nil {
		return t, fmt.Errorf("departure date: %w", err)
	}

	if t.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return t, fmt.Errorf("birth date: %w", err)
	}

	if t.PlannedDepartureDate, err = ParseSheetDate(get("Планова дата виїзду")); err != nil {
		return t, fmt.Errorf("planned departure date: %w", err)
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
		return t, fmt.Errorf("passport issue date: %w", err)
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
	t.Disability = ParseYesNo(get("Інвалідність (так/ні)"))
	t.Category = NullString(get("Категорія/коротка інформація про осіб"))
	t.FurtherAccommodation = NullString(get("Подальше розміщення"))
	t.Needs = NullString(get("Потреби"))

	return t, nil
}

func MainShelterMapper(object map[string]string) (MainShelterRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	var m MainShelterRecord
	var err error

	m.AddedBy = get("Внес дані")
	m.VPOStatus = get("Статус ВПО")
	m.LastName = get("Прізвище")
	m.FirstName = get("Ім'я")

	if m.Sex, err = ParseSex(get("Стать")); err != nil {
		return m, fmt.Errorf("sex: %w", err)
	}

	m.Category = get("Категорія")
	m.Status = get("Статус")

	m.IsTransit = ParseYesNo(get("Транзит"))

	if m.ArrivalDate, err = ParseSheetDate(get("Дата прибуття")); err != nil {
		return m, fmt.Errorf("arrival date: %w", err)
	}

	if m.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return m, fmt.Errorf("birth date: %w", err)
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
	m.Disability = ParseYesNo(get("Інвалідність"))
	m.EvacuatedFromTown = NullString(get("Населений пункт звідки евакуювали"))
	m.EvacuatedFromRegion = NullString(get("Область звідки евакуювали"))
	m.FurtherAccommodation = NullString(get("Подальше розміщення"))
	m.TaxID = NullString(get("ІПН"))
	m.Phone = NullString(get("Телефон"))
	m.FamilyComposition = NullString(get("Склад сім'ї"))
	m.RegistrationAddress = NullString(get("Адреса реєстрації"))
	m.Needs = NullString(get("Потреби"))

	if m.DepartureDate, err = ParseSheetDate(get("Дата виїзду")); err != nil {
		return m, fmt.Errorf("departure date: %w", err)
	}

	m.Passport.Number = NullString(get("Серія та номер паспорту"))
	m.Passport.Issuer = NullString(get("Ким виданий"))
	if m.Passport.IssueDate, err = ParseSheetDate(get("Коли виданий")); err != nil {
		return m, fmt.Errorf("passport issued date: %w", err)
	}

	return m, nil
}

func LocalTransportationsMapper(object map[string]string) (LocalTransportationRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	var l LocalTransportationRecord
	var err error

	l.LastName = get("Прізвище")
	l.FirstName = get("Ім'я")

	if l.Sex, err = ParseSex(get("Стать")); err != nil {
		return l, fmt.Errorf("sex: %w", err)
	}

	l.InformedConsent = ParseYesNo(get("Надаю інформо-вану згоду (так/ні)"))
	l.Disability = ParseYesNo(get("Інвалідність (так/ні)"))

	if l.Date, err = ParseSheetDate(get("Дата [ДД/ММ/РРРР]")); err != nil {
		return l, fmt.Errorf("transportation date: %w", err)
	}

	if l.BirthDate, err = ParseSheetDate(get("Дата народження [ДД/ММ/РРРР]")); err != nil {
		return l, fmt.Errorf("birth date: %w", err)
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

	return l, nil
}

func HotlineMapper(object map[string]string) (HotlineCallRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	parseInt := func(key string) int {
		val, _ := strconv.Atoi(get(key))
		return val
	}
	var h HotlineCallRecord
	var err error

	h.Responsible = get("Відповідальний")
	h.BeneficiaryFullName = NullString(get("ПІБ людини, якій буде надана послуга"))
	h.Gender, err = ParseSex(get("Стать"))
	if err != nil {
		return h, fmt.Errorf("sex: %w", err)
	}

	h.Service = get("Послуга")

	if h.Timestamp, err = ParseSheetDate(get("Позначка часу")); err != nil {
		return h, fmt.Errorf("timestamp: %w", err)
	}

	if h.CallDate, err = ParseSheetDate(get("Дата дзвінка")); err != nil {
		return h, fmt.Errorf("call date: %w", err)
	}

	if h.BirthDate, err = ParseSheetDate(get("Дата народження")); err != nil {
		return h, fmt.Errorf("birth date: %w", err)
	}

	if h.EvacuationDate, err = ParseSheetDate(get("Дата евакуації")); err != nil {
		return h, fmt.Errorf("evacuation date: %w", err)
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

	return h, nil
}

func VehicleTripMapper(object map[string]string) (VehicleTripRecord, error) {
	get := func(key string) string { return strings.TrimSpace(object[key]) }
	var v VehicleTripRecord
	parseOptionalInt := func(val string) int {
		clean := strings.TrimSpace(val)
		if clean == "" {
			return 0
		}
		num, err := strconv.Atoi(clean)
		if err != nil {
			return 0
		}
		return num
	}

	parseOptionalFloat := func(val string) float64 {
		clean := strings.TrimSpace(strings.ReplaceAll(val, ",", "."))
		if clean == "" {
			return 0.0
		}
		f, err := strconv.ParseFloat(clean, 64)
		if err != nil {
			return 0.0
		}
		return f
	}

	var err error

	//tsRaw := get("Позначка часу")
	//if tsRaw != "" {
	//
	//	if parsedTime, parseErr := time.Parse("02.01.2006 15:04:05", tsRaw); parseErr == nil {
	//		v.Timestamp = parsedTime
	//	} else if parsedTime, parseErr = time.Parse("2006-01-02 15:04:05", tsRaw); parseErr == nil {
	//		v.Timestamp = parsedTime
	//	}
	//}

	if v.Date, err = ParseSheetDate(get("Дата")); err != nil {
		return v, fmt.Errorf("trip date: %w", err)
	}

	v.StartTime = NullString(get("СТАРТ"))
	v.EndTime = NullString(get("ФІНІШ"))

	v.Driver = get("Водій/відповідальний за виїзд")
	v.Team = NullString(get("Команда виїзду"))
	v.VehiclePlate = get("Марка машини – номери")
	v.Route = get("Маршрут")
	v.TripType = NullString(get("Тип виїзду"))

	v.EvacueesCount = parseOptionalInt(get("Кількість евакуйованих осіб"))
	v.LowMobilityCount = parseOptionalInt(get("Кількість інвалідів І маломобільних"))
	v.ChildrenCount = parseOptionalInt(get("Кількість дітей"))

	v.OdometerStart = parseOptionalFloat(get("Показання oдометра на початку поїздки, км"))
	if v.OdometerStart == 0 {
		v.OdometerStart = parseOptionalFloat(get("Показання oдометра на\u00a0початкупоїздки, км"))
	}

	v.OdometerEnd = parseOptionalFloat(get("Показання oдометра на кінець поїздки, км"))
	if v.OdometerEnd == 0 {
		v.OdometerEnd = parseOptionalFloat(get("Показання oдометра на\u00a0кiнецьпоїздки, км"))
	}

	v.FuelRefueled = parseOptionalFloat(get("Заправлено палива, л"))
	v.FuelRemaining = parseOptionalFloat(get("Залишок палива(л.)"))

	v.Comments = NullString(get("Коментарі"))
	v.ExtraCol = NullString(get("Стовпець 16"))

	return v, nil
}
