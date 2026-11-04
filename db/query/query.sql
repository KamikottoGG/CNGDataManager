-- ============================================================================
-- 1. АНАЛИТИКА И СТАТИСТИКА ЭВАКУАЦИИ ПО ДНЯМ И ПЕРИОДАМ
-- ============================================================================

-- name: GetEvacuationDailyStats :many
-- Сводка по дням за выбранный диапазон дат: общее число, дети (до 18), люди с инвалидностью, разбивка по полу
SELECT
    e.evacuation_date,
    COUNT(e.id) AS total_evacuated,
    COUNT(DISTINCT e.beneficiary_id) AS unique_beneficiaries,
    CAST(COALESCE(SUM(CASE WHEN b.sex = 'FEMALE' THEN 1 ELSE 0 END), 0) AS SIGNED) AS female_count,
    CAST(COALESCE(SUM(CASE WHEN b.sex = 'MALE' THEN 1 ELSE 0 END), 0) AS SIGNED) AS male_count,
    CAST(COALESCE(SUM(CASE WHEN b.disability = 1 THEN 1 ELSE 0 END), 0) AS SIGNED) AS disability_count,
    CAST(COALESCE(SUM(CASE WHEN TIMESTAMPDIFF(YEAR, b.birth_date, e.evacuation_date) < 18 THEN 1 ELSE 0 END), 0) AS SIGNED) AS children_count,
    CAST(COALESCE(SUM(CASE WHEN TIMESTAMPDIFF(YEAR, b.birth_date, e.evacuation_date) >= 60 THEN 1 ELSE 0 END), 0) AS SIGNED) AS elderly_count
FROM evacuation_indicator e
         JOIN beneficiary b ON e.beneficiary_id = b.id
WHERE e.evacuation_date BETWEEN ? AND ?
GROUP BY e.evacuation_date
ORDER BY e.evacuation_date DESC;

-- name: GetEvacuationsCountByDayForSettlement :many
-- Количество эвакуированных посуточно из конкретного населенного пункта
SELECT
    e.evacuation_date,
    COUNT(e.id) AS count
FROM evacuation_indicator e
WHERE e.evacuated_from_settlement_id = ?
  AND e.evacuation_date BETWEEN ? AND ?
GROUP BY e.evacuation_date
ORDER BY e.evacuation_date ASC;

-- name: GetTopDestinationsByDateRange :many
-- Топ направлений эвакуации (куда вывозят) за определенный период
SELECT
    COALESCE(e.evacuated_to, 'Не указано') AS destination,
    COUNT(e.id) AS total_evacuated
FROM evacuation_indicator e
WHERE e.evacuation_date BETWEEN ? AND ?
GROUP BY e.evacuated_to
ORDER BY total_evacuated DESC
LIMIT ?;

-- ============================================================================
-- 2. ГЕОГРАФИЧЕСКАЯ АГРЕГАЦИЯ (РЕГИОНЫ, ГРОМАДЫ, НАСЕЛЕННЫЕ ПУНКТЫ)
-- ============================================================================

-- name: GetEvacuationStatsByRegion :many
-- Количество эвакуированных с группировкой по областям/регионам исхода
SELECT
    r.id AS region_id,
    r.name AS region_name,
    COUNT(e.id) AS total_evacuated
FROM evacuation_indicator e
         JOIN settlements s ON e.evacuated_from_settlement_id = s.id
         JOIN communities c ON s.community_id = c.id
         JOIN districts d ON c.district_id = d.id
         JOIN regions r ON d.region_id = r.id
WHERE e.evacuation_date BETWEEN ? AND ?
GROUP BY r.id, r.name
ORDER BY total_evacuated DESC;

-- name: GetEvacuationStatsByCommunity :many
-- Детализация по громадам внутри выбранного района
SELECT
    c.id AS community_id,
    c.name AS community_name,
    COUNT(e.id) AS total_evacuated
FROM evacuation_indicator e
         JOIN settlements s ON e.evacuated_from_settlement_id = s.id
         JOIN communities c ON s.community_id = c.id
WHERE c.district_id = ?
  AND e.evacuation_date BETWEEN ? AND ?
GROUP BY c.id, c.name
ORDER BY total_evacuated DESC;

-- ============================================================================
-- 3. ПОЛУЧЕНИЕ ДЕТАЛЬНЫХ ДАННЫХ И РЕЕСТРОВ (С JOIN)
-- ============================================================================

-- name: GetFullEvacuationRecord :one
-- Полная карточка конкретной записи об эвакуации со всей иерархией локаций
SELECT
    e.id AS evacuation_id,
    e.evacuation_date,
    e.evacuated_to,
    e.submitted_time,
    -- Бенефициар
    b.id AS beneficiary_id,
    b.name,
    b.surname,
    b.patronymic,
    b.phone,
    b.tax_id,
    b.passport_data,
    b.sex,
    b.birth_date,
    b.disability,
    -- Откуда эвакуирован
    from_s.id AS from_settlement_id,
    from_s.name AS from_settlement_name,
    from_s.type AS from_settlement_type,
    from_c.name AS from_community_name,
    from_d.name AS from_district_name,
    from_r.name AS from_region_name
FROM evacuation_indicator e
         JOIN beneficiary b ON e.beneficiary_id = b.id
         LEFT JOIN settlements from_s ON e.evacuated_from_settlement_id = from_s.id
         LEFT JOIN communities from_c ON from_s.community_id = from_c.id
         LEFT JOIN districts from_d ON from_c.district_id = from_d.id
         LEFT JOIN regions from_r ON from_d.region_id = from_r.id
WHERE e.id = ? LIMIT 1;

-- name: ListEvacuationsByDayDetailed :many
-- Список эвакуированных за конкретный день с постраничной пагинацией
SELECT
    e.id AS evacuation_id,
    e.evacuation_date,
    e.evacuated_to,
    b.id AS beneficiary_id,
    b.name,
    b.surname,
    b.phone,
    b.disability,
    s.name AS from_settlement_name
FROM evacuation_indicator e
         JOIN beneficiary b ON e.beneficiary_id = b.id
         LEFT JOIN settlements s ON e.evacuated_from_settlement_id = s.id
WHERE e.evacuation_date = ?
ORDER BY e.submitted_time DESC
LIMIT ? OFFSET ?;

-- name: GetBeneficiaryEvacuationHistory :many
-- История перемещений/эвакуаций одного конкретного человека
SELECT
    e.id,
    e.evacuation_date,
    e.evacuated_to,
    e.submitted_time,
    s.name AS from_settlement
FROM evacuation_indicator e
         LEFT JOIN settlements s ON e.evacuated_from_settlement_id = s.id
WHERE e.beneficiary_id = ?
ORDER BY e.evacuation_date DESC;

-- ============================================================================
-- 4. УПРАВЛЕНИЕ БЕНЕФИЦИАРАМИ И ОПЕРАЦИИ ВСТАВКИ
-- ============================================================================

-- name: CreateBeneficiary :execlastid
INSERT INTO beneficiary (
    name, surname, patronymic, registration_date, tax_id,
    passport_data, disability, sex, birth_date, phone,
    settlement_id, passport_issued_by, passport_issue_date, passport_issue_country
) VALUES (
             ?, ?, ?, ?, ?,
             ?, ?, ?, ?, ?,
             ?, ?, ?, ?
         );

-- name: CreateEvacuationIndicator :execlastid
INSERT INTO evacuation_indicator (
    evacuation_date, evacuated_to, beneficiary_id,
    submitted_time, evacuated_from_settlement_id
) VALUES (
             ?, ?, ?, ?, ?
         );

-- name: FindBeneficiaryByTaxOrPassport :one
-- Поиск дубликатов бенефициара перед регистрацией
SELECT id, name, surname, phone
FROM beneficiary
WHERE (tax_id IS NOT NULL AND tax_id = ?)
   OR (passport_data IS NOT NULL AND passport_data = ?)
LIMIT 1;

-- ============================================================================
-- 5. СПРАВОЧНИКИ (ДЛЯ ВЫПАДАЮЩИХ СПИСКОВ В UI)
-- ============================================================================

-- name: ListSettlementsByCommunity :many
SELECT id, name, type
FROM settlements
WHERE community_id = ?
ORDER BY name ASC;

-- name: GetSettlementHierarchy :one
-- Быстрое получение полной цепочки локации по settlement_id
SELECT
    s.id AS settlement_id,
    s.name AS settlement_name,
    s.type AS settlement_type,
    c.id AS community_id,
    c.name AS community_name,
    d.id AS district_id,
    d.name AS district_name,
    r.id AS region_id,
    r.name AS region_name
FROM settlements s
         JOIN communities c ON s.community_id = c.id
         JOIN districts d ON c.district_id = d.id
         JOIN regions r ON d.region_id = r.id
WHERE s.id = ? LIMIT 1;

-- ============================================================================
-- 6. ЧАТ (FUNKIN_CHAT_MESSAGES)
-- ============================================================================

-- name: CreateChatMessage :execlastid
INSERT INTO funkin_chat_messages (
    role, content, telegram_username, telegram_first_name, created_at
) VALUES (
             ?, ?, ?, ?, ?
         );

-- name: ListRecentChatMessages :many
-- Получение последних N сообщений для контекста
SELECT id, role, content, telegram_username, telegram_first_name, created_at
FROM funkin_chat_messages
ORDER BY created_at DESC
LIMIT ?;