-- +goose Up
-- +goose StatementBegin

create table branch
(
    branch_id  integer      not null,
    manager_id integer,
    name       varchar(128) not null,
    constraint branch_pk
        primary key (branch_id),
    constraint branch_manager_null_fk
        foreign key (manager_id) references manager
            on update cascade on delete set null
);

comment on table branch is 'Таблица, отвечающая за сущность филиала';

comment on column branch.branch_id is 'Идентификатор';

comment on column branch.manager_id is 'Идентификатор менеджера филиала';

comment on column branch.name is 'Название филиала';

alter table branch
    owner to postgres;

alter table branch
    add constraint check_name
        check (name ~~* '_%'::text);

comment on constraint check_name on branch is 'Проверка на непустоту названия';

alter table branch
    add constraint greater_than_zero
        check ((branch_id > 0) AND (manager_id > 0));

comment on constraint greater_than_zero on branch is 'Проверка на неотрицательность';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists branch
-- +goose StatementEnd
