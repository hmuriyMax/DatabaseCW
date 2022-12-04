-- +goose Up
-- +goose StatementBegin

create table device
(
    device_id      integer      not null,
    department_id  integer,
    name           varchar(128) not null,
    device_type_id integer      not null,
    constraint device_pk
        primary key (device_id),
    constraint device_department_null_fk
        foreign key (department_id) references department
            on update cascade on delete set null,
    constraint device_device_type_null_fk
        foreign key (device_type_id) references device_type
            on update cascade on delete cascade
);

comment on table device is 'Таблица с техникой';

comment on column device.device_id is 'Идентификатор';

comment on column device.department_id is 'Идентификатор отдела';

comment on column device.name is 'Название устройства';

alter table device
    owner to postgres;

alter table device
    add constraint greater_than_null
        check ((device_id > 0) AND (department_id > 0) AND (device_type_id > 0));

comment on constraint greater_than_null on device is 'Проверка на неотрицательность';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists device
-- +goose StatementEnd
