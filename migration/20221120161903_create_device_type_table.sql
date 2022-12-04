-- +goose Up
-- +goose StatementBegin

create table device_type
(
    device_type_id integer     not null,
    name           varchar(64) not null,
    constraint device_type_pk
        primary key (device_type_id)
);

comment on table device_type is 'Таблица для типа техники';

alter table device_type
    owner to postgres;

alter table device_type
    add constraint greater_than_zero
        check (device_type_id > 0);

comment on constraint greater_than_zero on device_type is 'Проверка на неотрицательность';

alter table device_type
    add constraint check_name
        check ((name)::text ~~* '_%'::text);

comment on constraint check_name on device_type is 'Имя не должно быть пустой строкой';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists device_type
-- +goose StatementEnd
