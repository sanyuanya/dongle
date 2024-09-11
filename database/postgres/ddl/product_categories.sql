-- auto-generated definition
create table product_categories
(
    snowflake_id varchar   default ''::character varying not null,
    name         varchar   default ''::character varying not null,
    status       integer   default 0                     not null,
    sorting      integer   default 0                     not null,
    created_at   timestamp default CURRENT_TIMESTAMP     not null,
    updated_at   timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at   timestamp
);

alter table product_categories
    owner to postgres;

