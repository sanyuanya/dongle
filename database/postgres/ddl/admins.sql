-- auto-generated definition
create table admins
(
    account      varchar  default ''::character varying           not null,
    password     varchar  default ''::character varying           not null,
    nick         varchar  default '系统管理员'::character varying not null,
    api_token    varchar  default ''::character varying           not null,
    created_at   timestamp                                        not null,
    updated_at   timestamp                                        not null,
    deleted_at   timestamp,
    snowflake_id varchar  default ''::character varying           not null,
    is_hidden    smallint default 0                               not null
);

comment on table admins is '系统管理员';

comment on column admins.account is '账号';

comment on column admins.password is '密码';

comment on column admins.nick is '昵称';

comment on column admins.api_token is 'API token';

comment on column admins.snowflake_id is '雪花ID';

alter table admins
    owner to postgres;

