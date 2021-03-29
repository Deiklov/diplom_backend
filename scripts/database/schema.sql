create table if not exists companies
(
    id char(36) not null
        constraint companies_pk
            primary key,
    name varchar(255) not null,
    description varchar default ''::character varying not null,
    country varchar(100) default ''::character varying not null,
    attributes json,
    ipo date,
    ticker char(5),
    logo varchar(255),
    weburl varchar(255)
);

comment on table companies is '"currency+industry+exchange to json attributes"';

alter table companies owner to andrey;

create unique index if not exists companies_id_uindex
    on companies (id);

create unique index if not exists companies_name_uindex
    on companies (name);

