-- auto-generated definition
create table pay
(
    snowflake_id   varchar default ''::character varying not null,
    name           varchar default ''::character varying not null,
    remark         varchar default ''::character varying not null,
    total_amount   bigint  default 0                     not null,
    total_num      integer default 0                     not null,
    status         varchar                               not null,
    created_at     timestamp                             not null,
    updated_at     timestamp                             not null,
    deleted_at     timestamp,
    batch_id       varchar default ''::character varying not null,
    close_reason   varchar default ''::character varying not null,
    success_amount bigint  default 0                     not null,
    success_num    integer default 0                     not null,
    fail_amount    bigint  default 0                     not null,
    fail_num       integer default 0                     not null,
    create_time    varchar default ''::character varying not null,
    update_time    varchar default ''::character varying not null
);

comment on table pay is '支付信息';

comment on column pay.snowflake_id is '批次ID';

comment on column pay.name is '批次名称';

comment on column pay.remark is '批次备注';

comment on column pay.total_amount is '转账总金额';

comment on column pay.total_num is '转账总笔数';

comment on column pay.status is '批次状态';

comment on column pay.batch_id is '微信批次单号';

comment on column pay.close_reason is '批次关闭原因';

comment on column pay.success_amount is '转账成功金额';

comment on column pay.success_num is '转账成功笔数';

comment on column pay.fail_amount is '转账失败金额';

comment on column pay.fail_num is '转账失败笔数';

alter table pay
    owner to postgres;

