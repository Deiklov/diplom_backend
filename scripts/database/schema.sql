create table users
(
    id char(36) not null
        constraint user_pk
            primary key,
    email varchar(150) not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP,
    name varchar(100) default ''::character varying not null,
    password varchar(255) not null
);

comment on table users is 'Таблица пользователей';

alter table users owner to andrey;

create unique index user_id_uindex
    on users (id);

create unique index users_email_uindex
    on users (email);

create table companies
(
    id char(36) not null
        constraint companies_pk
            primary key,
    name varchar(255) not null,
    year smallint,
    description varchar default ''::character varying not null
);

comment on table companies is 'Компании';

alter table companies owner to andrey;

create unique index companies_id_uindex
    on companies (id);

create unique index companies_name_uindex
    on companies (name);

create table profiles
(
    id char(36) not null,
    user_id char(36) not null
        constraint profile_pk
            primary key
        constraint profile_users_id_fk
            references users
            on update cascade on delete cascade,
    age smallint,
    avatar_path varchar(400),
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP,
    deleted_at timestamp default CURRENT_TIMESTAMP
);

comment on table profiles is 'Профиль к юзеру 1 к 1';

alter table profiles owner to andrey;

create unique index profile_id_uindex
    on profiles (id);

create unique index profile_user_id_uindex
    on profiles (user_id);

create table company_by_users
(
    id char(36) not null
        constraint company_by_users_pk
            primary key,
    company_id char(36) not null
        constraint company_by_users_companies_id_fk
            references companies
            on update cascade on delete cascade,
    user_id char(36) not null
        constraint company_by_users_users_id_fk
            references users
            on update cascade on delete cascade
);

comment on table company_by_users is 'Избранные компании для юзеров';

alter table company_by_users owner to andrey;

create unique index company_by_users_id_uindex
    on company_by_users (id);

