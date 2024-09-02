-- auto-generated definition
create table product
(
    snowflake_id varchar default ''::character varying not null,
    name         varchar default ''::character varying not null,
    integral     bigint  default 0                     not null,
    created_at   timestamp                             not null,
    updated_at   timestamp                             not null,
    deleted_at   timestamp
);

comment on table product is '产品';

alter table product
    owner to postgres;

