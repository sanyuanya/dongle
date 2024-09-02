-- auto-generated definition
create table permissions
(
    snowflake_id varchar default ''::character varying not null,
    summary      varchar default ''::character varying not null,
    path         varchar default ''::character varying not null,
    created_at   timestamp                             not null,
    updated_at   timestamp                             not null,
    deleted_at   varchar,
    api_path     varchar default ''::character varying not null
);

comment on table permissions is '权限';

alter table permissions
    owner to postgres;

