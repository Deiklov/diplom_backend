create table users
(
    id char(36) not null
        constraint user_pk
            primary key,
    phone varchar(50) not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    updated_at timestamp default CURRENT_TIMESTAMP,
    name varchar(100) default ''::character varying not null
);

comment on table users is 'Таблица пользователей';

alter table users owner to andrey;

create unique index user_id_uindex
    on users (id);

