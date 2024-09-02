-- auto-generated definition
create table roles
(
    snowflake_id varchar  default ''::character varying not null,
    name         varchar  default ''::character varying not null,
    created_at   timestamp                              not null,
    updated_at   timestamp                              not null,
    deleted_at   timestamp,
    is_hidden    smallint default 0                     not null
);

comment on table roles is '角色';

comment on column roles.name is '角色名称';

alter table roles
    owner to postgres;

