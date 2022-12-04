-- +goose Up
-- +goose StatementBegin

create table performer
(
    worker_id integer      not null
        constraint performer_pk
            primary key,
    name      varchar(128) not null,
    constraint check_name
        check (name ilike '_% _% _%')
);

comment on table performer is 'Таблица работника';

comment on column performer.worker_id is 'Идентификатор';

comment on column performer.name is 'Имя работника';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists performer;
-- +goose StatementEnd
