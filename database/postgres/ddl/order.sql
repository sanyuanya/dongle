-- auto-generated definition
create table "order"
(
    snowflake_id     varchar   default ''::character varying not null,
    transaction_id   varchar   default ''::character varying not null,
    app_id           varchar   default ''::character varying not null,
    mch_id           varchar   default ''::character varying not null,
    trade_type       varchar   default ''::character varying not null,
    trade_state      varchar   default ''::character varying not null,
    trade_state_desc varchar   default ''::character varying not null,
    bank_type        varchar   default ''::character varying not null,
    success_time     varchar   default ''::character varying not null,
    open_id          varchar   default ''::character varying not null,
    user_id          varchar   default ''::character varying not null,
    total            bigint    default 0                     not null,
    payer_total      bigint    default 0                     not null,
    currency         varchar   default ''::character varying not null,
    payer_currency   varchar   default ''::character varying not null,
    out_trade_no     varchar   default ''::character varying not null,
    created_at       timestamp default CURRENT_TIMESTAMP     not null,
    updated_at       timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at       timestamp,
    prepay_id        varchar   default ''::character varying not null,
    expiration_time  integer   default 0                     not null,
    address_id       varchar   default ''::character varying not null,
    consignee        varchar   default ''::character varying not null,
    phone_number     varchar   default ''::character varying not null,
    location         varchar   default ''::character varying not null,
    detailed_address varchar   default ''::character varying not null,
    order_state      integer   default 0                     not null
);

comment on table "order" is '订单';

comment on column "order".transaction_id is '微信支付系统生成的订单号';

comment on column "order".trade_type is '交易类型';

comment on column "order".trade_state is '交易状态';

comment on column "order".trade_state_desc is '交易状态描述';

comment on column "order".bank_type is '银行类型';

comment on column "order".success_time is '支付完成时间';

comment on column "order".open_id is '支付者信息';

comment on column "order".user_id is '用户id';

comment on column "order".total is '订单总金额，单位为分';

comment on column "order".payer_total is '用户支付金额，单位为分';

comment on column "order".currency is 'CNY：人民币';

comment on column "order".payer_currency is '用户支付币种';

comment on column "order".out_trade_no is '商户订单号';

comment on column "order".prepay_id is '预支付交易会话标识';

comment on column "order".address_id is '收获地址';

comment on column "order".consignee is '收货人';

comment on column "order".phone_number is '收货人手机号';

comment on column "order".location is '收货人所在地区';

comment on column "order".detailed_address is '收货人-详细地址';

comment on column "order".order_state is '订单状态';

alter table "order"
    owner to postgres;

