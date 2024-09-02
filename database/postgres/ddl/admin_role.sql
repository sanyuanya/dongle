-- auto-generated definition
create table admin_role
(
    snowflake_id varchar default ''::character varying not null,
    admin_id     varchar default ''::character varying not null,
    role_id      varchar default ''::character varying not null,
    created_at   timestamp                             not null,
    updated_at   timestamp                             not null,
    deleted_at   timestamp
);

comment on table admin_role is '账号角色关联';

alter table admin_role
    owner to postgres;

