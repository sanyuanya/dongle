-- auto-generated definition
create table role_permission
(
    snowflake_id  varchar default ''::character varying not null,
    role_id       varchar default ''::character varying not null,
    permission_id varchar default ''::character varying not null,
    created_at    timestamp                             not null,
    updated_at    timestamp                             not null,
    deleted_at    timestamp
);

comment on table role_permission is '角色权限表';

alter table role_permission
    owner to postgres;

