-- +goose Up
-- +goose StatementBegin

create table manager
(
    manager_id integer      not null
        constraint id
            primary key
        constraint greater_than_zero
            check (manager_id > 0),
    name       varchar(128) not null
        constraint check_name
            check ((name)::text ~~* '_% _% _%'::text)
);

comment on table manager is 'Менеджер';

comment on column manager.manager_id is 'Идентификатор';

comment on column manager.name is 'Имя руководителя';

alter table manager
    owner to postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists manager
-- +goose StatementEnd
