create table platform.user_equipment_info
(
    id                  bigserial           not null primary key,
    base_id             bigint              not null default 0,
    type                integer             not null default 0,
    position            integer             not null default 0,
    quality             integer             not null default 0,
    user_id          	bigint              not null default 0,
    user_hero_id        bigint              not null default 0,
    hero_id         	bigint              not null default 0,
    level               integer             not null default 0,
    star                integer             not null default 0,
    status              smallint            not null default 0,
    drop_attrs          json,
    upgrade_attrs       json,
    create_time         timestamp without time zone not null default CURRENT_TIMESTAMP,
    update_time         timestamp without time zone not null default CURRENT_TIMESTAMP
);

create index userequipmentinfo_user_id on platform.user_equipment_info (user_id);
create index userequipmentinfo_base_id on platform.user_equipment_info (base_id);
create index userequipmentinfo_user_hero_id on platform.user_equipment_info (user_hero_id);
create index userequipmentinfo_hero_id on platform.user_equipment_info (hero_id);
create index userequipmentinfo_create_time on platform.user_equipment_info (create_time);
create index userequipmentinfo_update_time on platform.user_equipment_info (update_time);