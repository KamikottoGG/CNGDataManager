package sheets

var (
	_ RecordMapper[EvacuationRecord]          = EvacuationMapper
	_ RecordMapper[TransitHubRecord]          = TransitHubMapper
	_ RecordMapper[MainShelterRecord]         = MainShelterMapper
	_ RecordMapper[LocalTransportationRecord] = LocalTransportationsMapper
	_ RecordMapper[HotlineCallRecord]         = HotlineMapper
	_ RecordMapper[VehicleTripRecord]         = VehicleTripMapper
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

type MainShelterRecord struct {
	AddedBy              string         `json:"added_by"`              // Внес дані
	ArrivalDate          DateOnly       `json:"arrival_date"`          // Дата прибуття (формат YYYY-MM-DD или DD.MM.YYYY)
	VPOStatus            string         `json:"vpo_status"`            // Статус ВПО
	IsTransit            bool           `json:"is_transit"`            // Транзит (true/false)
	LastName             string         `json:"last_name"`             // Прізвище
	FirstName            string         `json:"first_name"`            // Ім'я
	MiddleName           NullString     `json:"middle_name"`           // По батькові (может отсутствовать -> null)
	BirthDate            DateOnly       `json:"birth_date"`            // Дата народження
	Age                  int            `json:"age"`                   // Вік
	Sex                  BeneficiarySex `json:"sex"`                   // Стать
	Category             string         `json:"category"`              // Категорія
	Disability           bool           `json:"disability"`            // Інвалідність (указатель, если null/нет)
	EvacuatedFromTown    NullString     `json:"evacuated_from_town"`   // Населений пункт звідки евакуювали
	EvacuatedFromRegion  NullString     `json:"evacuated_from_region"` // Область звідки евакуювали
	FurtherAccommodation NullString     `json:"further_accommodation"` // Подальше розміщення
	TaxID                NullString     `json:"tax_id"`                // ІПН (может быть пустым или содержать буквы для отказников)
	Phone                NullString     `json:"phone"`                 // Телефон
	Passport             BeneficiaryPassport
	FamilyComposition    NullString `json:"family_composition"`   // Склад сім'ї
	RegistrationAddress  NullString `json:"registration_address"` // Адреса реєстрації
	Needs                NullString `json:"needs"`                // Потреби
	Status               string     `json:"status"`               // Статус (напр. "Проживает", "Выехал")
	DepartureDate        DateOnly   `json:"departure_date"`       // Дата виїзду (если еще не выехал -> null)
}

type TransitHubRecord struct {
	ResidenceStatus        string              `json:"residence_status"`        // Статус проживання
	ArrivalDate            DateOnly            `json:"arrival_date"`            // Дата прибуття
	DepartureDate          DateOnly            `json:"departure_date"`          // Дата вибуття
	FromRegion             NullString          `json:"from_region"`             // Звідки область
	FromDistrict           NullString          `json:"from_district"`           // Звідки район
	FromCommunity          NullString          `json:"from_community"`          // Звідки громада
	FromSettlement         NullString          `json:"from_settlement"`         // Звідки нас пункт
	LastName               string              `json:"last_name"`               // Прізвище
	FirstName              string              `json:"first_name"`              // Ім'я
	MiddleName             NullString          `json:"middle_name"`             // По батькові
	Sex                    BeneficiarySex      `json:"sex"`                     // Стать
	BirthDate              DateOnly            `json:"birth_date"`              // Дата народження
	Age                    int                 `json:"age"`                     // Вік
	TaxID                  NullString          `json:"tax_id"`                  // ІПН
	Passport               BeneficiaryPassport `json:"passport"`                // Паспортні данные
	RegistrationRegion     NullString          `json:"registration_region"`     // Прописка: Область
	RegistrationCommunity  NullString          `json:"registration_community"`  // Прописка: Громада
	RegistrationSettlement NullString          `json:"registration_settlement"` // Прописка: Нас пункт
	RegistrationAddress    NullString          `json:"registration_address"`    // Прописка: Адреса
	Phone                  NullString          `json:"phone"`                   // Телефон
	Disability             bool                `json:"disability"`              // Інвалідність (так/ні)
	Category               NullString          `json:"category"`                // Категорія/коротка інформація про осіб
	FurtherAccommodation   NullString          `json:"further_accommodation"`   // Подальше розміщення
	PlannedDepartureDate   DateOnly            `json:"planned_departure_date"`  // Планова дата виїзду
	Needs                  NullString          `json:"needs"`                   // Потреби
}

type VehicleTripRecord struct {
	Date DateOnly `json:"date"`

	StartTime NullString `json:"startTime"` // СТАРТ
	EndTime   NullString `json:"endTime"`   // ФІНІШ

	Driver       string     `json:"driver"`       // Водій/відповідальний за виїзд
	Team         NullString `json:"team"`         // Команда виїзду
	VehiclePlate string     `json:"vehiclePlate"` // Марка машини – номери
	Route        string     `json:"route"`        // Маршрут
	TripType     NullString `json:"tripType"`     // Тип виїзду

	EvacueesCount    int `json:"evacueesCount"`    // Кількість евакуйованих осіб
	LowMobilityCount int `json:"lowMobilityCount"` // Кількість інвалідів І маломобільних
	ChildrenCount    int `json:"childrenCount"`    // Кількість дітей

	OdometerStart float64 `json:"odometerStart"` // Показання одометра на початку поїздки, км
	OdometerEnd   float64 `json:"odometerEnd"`   // Показання одометра на кінець поїздки, км
	FuelRefueled  float64 `json:"fuelRefueled"`  // Заправлено палива, л
	FuelRemaining float64 `json:"fuelRemaining"` // Залишок палива(л.)

	Comments NullString `json:"comments"`   // Коментарі
	ExtraCol NullString `json:"extraCol16"` // Стовпець 16 (технічна або запасна колонка)
}
