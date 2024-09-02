-- auto-generated definition
create table withdrawals
(
    snowflake_id      varchar default ''::character varying       not null,
    user_id           varchar default ''::character varying       not null,
    integral          integer default 0                           not null,
    withdrawal_method varchar default 'wechat'::character varying not null,
    life_cycle        integer default 0                           not null,
    created_at        timestamp                                   not null,
    updated_at        timestamp                                   not null,
    deleted_at        timestamp,
    rejection         varchar default ''::character varying       not null,
    detail_id         varchar default ''::character varying       not null,
    pay_id            varchar default ''::character varying       not null,
    initiate_time     varchar default ''::character varying       not null,
    update_time       varchar default ''::character varying       not null,
    open_id           varchar default ''::character varying       not null,
    mch_id            varchar default ''::character varying       not null,
    payment_status    varchar default ''::character varying       not null
);

comment on table withdrawals is '提现记录';

comment on column withdrawals.snowflake_id is '明细号';

comment on column withdrawals.user_id is '用户ID';

comment on column withdrawals.integral is '提现积分';

comment on column withdrawals.withdrawal_method is '提现方式';

comment on column withdrawals.life_cycle is '1 申请提现 2 驳回 3 审批通过 4 失败';

comment on column withdrawals.detail_id is '微信明细单号';

comment on column withdrawals.pay_id is '批次号';

comment on column withdrawals.initiate_time is '转账发起时间';

comment on column withdrawals.update_time is '明细更新时间';

alter table withdrawals
    owner to postgres;

