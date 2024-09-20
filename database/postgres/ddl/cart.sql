-- auto-generated definition
create table cart
(
    snowflake_id varchar   default ''::character varying not null,
    commodity_id varchar   default ''::character varying not null,
    sku_id       varchar   default ''::character varying not null,
    quantity     integer   default 1                     not null,
    user_id      varchar   default ''::character varying not null,
    created_at   timestamp default CURRENT_TIMESTAMP     not null,
    updated_at   timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at   timestamp
);

comment on table cart is '购物车';

comment on column cart.commodity_id is '商品ID';

comment on column cart.sku_id is 'sku id ';

comment on column cart.quantity is '数量';

comment on column cart.user_id is '用户ID';

alter table cart
    owner to postgres;

