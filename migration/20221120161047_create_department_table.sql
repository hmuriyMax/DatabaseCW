-- +goose Up
-- +goose StatementBegin

create table department
(
    department_id integer     not null
        constraint department_pk
            primary key,
    branch_id     integer
        constraint department_branch_null_fk
            references branch
            on update cascade on delete cascade,
    name          varchar(64) not null
        constraint check_name
            check ((name)::text ~~* '_%'::text),
    constraint greater_than_zero
        check ((department_id > 0) AND (branch_id > 0))
);

comment on table department is 'Отдел';

comment on column department.department_id is 'Идентификатор';

comment on column department.branch_id is 'Идентификатор филиала';

comment on column department.name is 'Название отдела';

comment on constraint greater_than_zero on department is 'Проверка на неотрицательность';

alter table department
    owner to postgres;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists department
-- +goose StatementEnd
