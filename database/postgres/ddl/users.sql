-- auto-generated definition
create table users
(
    nick                varchar  default ''::character varying not null,
    phone               varchar  default ''::character varying not null,
    openid              varchar  default ''::character varying not null,
    session_key         varchar  default ''::character varying not null,
    api_token           varchar  default ''::character varying not null,
    is_white            smallint default 0                     not null,
    avatar              varchar  default ''::character varying not null,
    integral            integer  default 0                     not null,
    shipments           integer  default 0                     not null,
    province            varchar  default ''::character varying not null,
    city                varchar  default ''::character varying not null,
    created_at          timestamp                              not null,
    updated_at          timestamp                              not null,
    deleted_at          timestamp,
    district            varchar  default ''::character varying not null,
    id_card             varchar  default ''::character varying not null,
    company_name        varchar  default ''::character varying not null,
    job                 varchar  default ''::character varying not null,
    alipay_account      varchar  default ''::character varying not null,
    snowflake_id        varchar  default ''::character varying not null,
    withdrawable_points integer  default 0                     not null
);

comment on table users is '用户信息';

comment on column users.nick is '昵称';

comment on column users.phone is '手机号';

comment on column users.api_token is '接口token';

comment on column users.is_white is '是否是白名单';

comment on column users.avatar is '头像';

comment on column users.integral is '积分';

comment on column users.shipments is '出货量';

comment on column users.province is '省';

comment on column users.city is '城市';

comment on column users.district is '区域';

comment on column users.id_card is '身份证';

comment on column users.company_name is '单位名称';

comment on column users.job is '职务';

comment on column users.alipay_account is '支付宝账号';

comment on column users.withdrawable_points is '可提现积分';

alter table users
    owner to postgres;

