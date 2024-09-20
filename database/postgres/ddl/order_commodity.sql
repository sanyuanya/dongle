-- auto-generated definition
create table order_commodity
(
    snowflake_id          varchar   default ''::character varying not null,
    commodity_id          varchar   default ''::character varying not null,
    commodity_name        varchar   default ''::character varying not null,
    commodity_code        varchar   default ''::character varying not null,
    categories_id         varchar   default ''::character varying not null,
    commodity_description varchar   default ''::character varying not null,
    sku_id                varchar   default ''::character varying not null,
    sku_code              varchar   default ''::character varying not null,
    price                 numeric   default 0                     not null,
    object_name           varchar   default ''::character varying not null,
    bucket_name           varchar   default ''::character varying not null,
    order_id              varchar   default ''::character varying not null,
    created_at            timestamp default CURRENT_TIMESTAMP     not null,
    updated_at            timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at            timestamp,
    quantity              integer   default 1                     not null,
    sku_name              varchar   default ''::character varying not null
);

comment on table order_commodity is '订单商品关联关系表';

comment on column order_commodity.commodity_id is '商品ID';

comment on column order_commodity.commodity_name is '商品名称';

comment on column order_commodity.commodity_code is '商品编号';

comment on column order_commodity.categories_id is '商品分类ID';

comment on column order_commodity.commodity_description is '商品描述';

comment on column order_commodity.sku_id is '商品SKU ID';

comment on column order_commodity.sku_code is '商品 SKU 编码';

comment on column order_commodity.price is '商品 sku 价格';

comment on column order_commodity.object_name is '对象名称';

comment on column order_commodity.bucket_name is '桶名称';

comment on column order_commodity.order_id is '订单ID';

comment on column order_commodity.quantity is '商品数量';

comment on column order_commodity.sku_name is 'sku名称';

alter table order_commodity
    owner to postgres;

