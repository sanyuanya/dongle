-- auto-generated definition
create table income_expense
(
    snowflake_id        varchar   default ''::character varying not null,
    summary             varchar   default ''::character varying not null,
    integral            integer   default 0                     not null,
    shipments           integer   default 0                     not null,
    user_id             varchar   default ''::character varying not null,
    created_at          timestamp                               not null,
    updated_at          timestamp                               not null,
    deleted_at          timestamp,
    batch               varchar   default ''::character varying not null,
    product_id          varchar   default ''::character varying not null,
    product_integral    bigint    default 0                     not null,
    importd_at          timestamp default CURRENT_TIMESTAMP     not null,
    withdrawable_points bigint    default 0                     not null,
    path                varchar   default ''::character varying not null,
    file_name           varchar   default ''::character varying not null
);

comment on table income_expense is '分红明细';

comment on column income_expense.snowflake_id is '雪花ID';

comment on column income_expense.summary is '描述';

comment on column income_expense.integral is '积分';

comment on column income_expense.shipments is '出货量';

comment on column income_expense.user_id is '用户ID';

comment on column income_expense.batch is '批次ID';

alter table income_expense
    owner to postgres;

