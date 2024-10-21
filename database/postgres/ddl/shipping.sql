-- auto-generated definition
create table shipping
(
    snowflake_id      varchar   default ''::character varying not null,
    order_id          varchar   default ''::character varying not null,
    task_id           varchar   default ''::character varying not null,
    third_order_id    varchar   default ''::character varying not null,
    order_number      varchar   default ''::character varying not null,
    e_order           text      default ''::text              not null,
    created_at        timestamp default CURRENT_TIMESTAMP     not null,
    updated_at        timestamp default CURRENT_TIMESTAMP     not null,
    deleted_at        timestamp,
    status            integer   default 0                     not null,
    user_cancel_msg   varchar   default ''::character varying not null,
    system_cancel_msg varchar   default ''::character varying not null,
    courier_name      varchar   default ''::character varying not null,
    courier_mobile    varchar   default ''::character varying not null,
    net_tel           varchar   default ''::character varying not null,
    net_code          varchar   default ''::character varying not null,
    weight            varchar   default ''::character varying not null,
    def_price         varchar   default ''::character varying not null,
    volume            varchar   default ''::character varying not null,
    actual_weight     varchar   default ''::character varying not null,
    print_task_id     varchar   default ''::character varying not null,
    label             varchar   default ''::character varying not null,
    pickup_code       varchar   default ''::character varying not null
);

comment on table shipping is '发货表';

comment on column shipping.order_id is '订单编号';

comment on column shipping.task_id is '快递100 任务ID';

comment on column shipping.third_order_id is '快递100 订单ID';

comment on column shipping.order_number is '快递单号';

comment on column shipping.e_order is '快递面单附属属性，根据各个快递公司返回属性';

comment on column shipping.status is '订单状态说明： 0：''下单成功''； 1：''已接单''； 2：''收件中''； 9：''用户主动取消''；10：''已取件''； 11：''揽货失败''；12：''已退回''； 13：''已签收''； 14：''异常签收''；15：''已结算'' ；99：''订单已取消''；101：''运输中''；200：''已出单''；201：''出单失败''；610：''下单失败''；155：''修改重量''(注意需要在工单系统中发起异常反馈并由快递100服务人员确认调重后才会有此状态回调，回调内容包含修改重量后的重量、运费、费用明细、业务类型)；166：订单复活（订单被取消，但是实际包裹已经发出，正常计费）；400：派送中';

comment on column shipping.user_cancel_msg is '用户取消原因';

comment on column shipping.system_cancel_msg is '系统取消或下单失败原因';

comment on column shipping.courier_name is '快递员姓名';

comment on column shipping.courier_mobile is '快递员电话';

alter table shipping
    owner to postgres;

