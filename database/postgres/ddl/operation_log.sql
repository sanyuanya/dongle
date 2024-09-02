-- auto-generated definition
create table operation_log
(
    snowflake_id              varchar default ''::character varying not null,
    operation_id              varchar default ''::character varying not null,
    income_expense_id         varchar default ''::character varying not null,
    user_id                   varchar default ''::character varying not null,
    before_updating_shipments integer default 0                     not null,
    after_updating_shipments  integer default 0                     not null,
    summary                   varchar default ''::character varying not null,
    created_at                timestamp                             not null,
    updated_at                timestamp                             not null,
    deleted_at                timestamp
);

comment on table operation_log is '操作日志';

comment on column operation_log.operation_id is '操作员ID';

comment on column operation_log.income_expense_id is '明细ID';

comment on column operation_log.user_id is '客户ID';

comment on column operation_log.before_updating_shipments is '更新前-出货量';

comment on column operation_log.after_updating_shipments is '更新后-出货量';

comment on column operation_log.summary is '概括';

alter table operation_log
    owner to postgres;

