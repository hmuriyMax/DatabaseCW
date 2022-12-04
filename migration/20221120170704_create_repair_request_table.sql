-- +goose Up
-- +goose StatementBegin

create table repair_request
(
    request_id     integer not null,
    device_id      integer not null,
    repair_type_id integer not null,
    worker_id      integer,
    duration       interval,
    constraint repair_request_pk
        primary key (request_id),
    constraint repair_request_device_null_fk
        foreign key (device_id) references device
            on update cascade on delete cascade,
    constraint repair_request_performer_null_fk
        foreign key (worker_id) references performer
            on update cascade on delete set null,
    constraint repair_request_repair_type_null_fk
        foreign key (repair_type_id) references repair_type
            on update cascade on delete cascade
);

comment on table repair_request is 'Запрос на ремонт';

comment on column repair_request.request_id is 'Идентификатор';

comment on column repair_request.device_id is 'Идентификатор устройства';

comment on column repair_request.repair_type_id is 'Идентификатор типа';

comment on column repair_request.worker_id is 'Идентификатор работиника';

comment on column repair_request.duration is 'Время окончания работ';

alter table repair_request
    owner to postgres;

alter table repair_request
    add constraint greater_than_zero
        check ((request_id > 0) AND (device_id > 0) AND (repair_type_id > 0) AND (worker_id > 0));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists repair_request;
-- +goose StatementEnd
