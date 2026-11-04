create table funkin_chat_messages
(
    id                  bigint auto_increment
        primary key,
    role                varchar(10)  not null,
    content             mediumtext   not null,
    telegram_username   varchar(255) null,
    telegram_first_name varchar(255) null,
    created_at          timestamp    not null
)
    collate = utf8mb4_unicode_ci;

create table regions
(
    id   bigint auto_increment
        primary key,
    name varchar(100) not null,
    constraint name
        unique (name)
);

create table districts
(
    id        bigint auto_increment
        primary key,
    name      varchar(100) not null,
    region_id bigint       not null,
    constraint districts_ibfk_1
        foreign key (region_id) references regions (id)
);

create table communities
(
    id          bigint auto_increment
        primary key,
    name        varchar(100) not null,
    district_id bigint       not null,
    constraint communities_ibfk_1
        foreign key (district_id) references districts (id)
);

create index district_id
    on communities (district_id);

create index region_id
    on districts (region_id);

create table settlements
(
    id           bigint auto_increment
        primary key,
    name         varchar(100)                     not null,
    type         enum ('CITY', 'TOWN', 'VILLAGE') not null,
    community_id bigint                           not null,
    constraint settlements_ibfk_1
        foreign key (community_id) references communities (id)
);

create table beneficiary
(
    id                     bigint auto_increment
        primary key,
    name                   varchar(255)            not null,
    surname                varchar(255)            not null,
    patronymic             varchar(255)            null,
    registration_date      date                    null,
    tax_id                 varchar(10)             null,
    passport_data          varchar(255)            null,
    disability             tinyint(1)              not null,
    sex                    enum ('MALE', 'FEMALE') not null,
    birth_date             date                    not null,
    phone                  varchar(30)             null,
    settlement_id          bigint                  null,
    passport_issued_by     varchar(255)            null,
    passport_issue_date    date                    null,
    passport_issue_country varchar(100)            null,
    constraint beneficiary_ibfk_1
        foreign key (settlement_id) references settlements (id)
);

create index settlement_id
    on beneficiary (settlement_id);

create table evacuation_indicator
(
    id                           bigint auto_increment
        primary key,
    evacuation_date              date         not null,
    evacuated_to                 varchar(100) null,
    beneficiary_id               bigint       not null,
    submitted_time               timestamp    not null,
    evacuated_from_settlement_id bigint       null,
    constraint evacuation_indicator_ibfk_1
        foreign key (beneficiary_id) references beneficiary (id),
    constraint fk_evac_from_settlement
        foreign key (evacuated_from_settlement_id) references settlements (id)
);

create index beneficiary_id
    on evacuation_indicator (beneficiary_id);

create index community_id
    on settlements (community_id);

