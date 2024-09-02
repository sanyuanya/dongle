-- auto-generated definition
create table mini
(
    app_id       varchar default ''::character varying not null,
    secret       varchar default ''::character varying not null,
    access_token varchar default ''::character varying not null,
    expires_in   integer default 0                     not null
);

comment on table mini is '配置信息';

comment on column mini.app_id is '微信 app id ';

comment on column mini.secret is '微信 secret';

comment on column mini.access_token is 'client 访问令牌';

comment on column mini.expires_in is '过期时间';

alter table mini
    owner to postgres;

