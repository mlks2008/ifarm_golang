create table app.gateway_user
(
    id              bigserial                   not null primary key,
    real_userid     varchar(100)                not null default '',
    token_userid    varchar(100)                not null default '',
    web3_address    varchar(100)                not null default '',
    hd_address      varchar(100)                not null default '',
    create_time     timestamp without time zone not null default CURRENT_TIMESTAMP,
    update_time     timestamp without time zone not null default CURRENT_TIMESTAMP
);

create unique index gwuser_real_userid on app.gateway_user (real_userid);
create unique index gwuser_token_userid on app.gateway_user (token_userid);
create index gwuser_web3_address on app.gateway_user (web3_address);
create index gwuser_hd_address on app.gateway_user (hd_address);
create index gwuser_create_time on app.gateway_user (create_time);
create index gwuser_update_time on app.gateway_user (update_time);