-- +goose Up
-- +goose StatementBegin

create table repair_type
(
    repair_type_id integer      not null
        constraint repair_type_pk
            primary key,
    duration       interval     not null,
    name           varchar(128) not null,
    constraint check_name
        check (name ilike '_%'),
    constraint greater_than_zero
        check (repair_type_id > 0)
);

comment on table repair_type is 'Тип ремонта';

comment on column repair_type.repair_type_id is 'Идентификатор';

comment on column repair_type.duration is 'Длительность ремонта для данного типа';

comment on column repair_type.name is 'Имя типа';

comment on constraint greater_than_zero on repair_type is 'Проверка на неотрицательность';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists repair_type;
-- +goose StatementEnd
