-- auto-generated definition
create table address
(
    snowflake_id     varchar   default ''::character varying not null,
    consignee        varchar   default ''::character varying not null,
    phone_number     varchar   default ''::character varying not null,
    location         varchar   default ''::character varying not null,
    detailed_address varchar   default ''::character varying not null,
    longitude        bigint    default 0                     not null,
    latitude         bigint    default 0                     not null,
    user_id          varchar   default ''::character varying not null,
    is_default       integer   default 0                     not null,
    created_at       timestamp default CURRENT_TIMESTAMP     not null,
    updated_at       timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at       timestamp
);

comment on table address is '收获地址';

comment on column address.consignee is '收货人';

comment on column address.phone_number is '手机号';

comment on column address.location is '所在地区';

comment on column address.longitude is '经度';

comment on column address.user_id is '用户ID';

alter table address
    owner to postgres;

