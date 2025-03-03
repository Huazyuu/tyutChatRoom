drop database go_gin_chat;
create database go_gin_chat;
use go_gin_chat;
set names utf8mb4;
set foreign_key_checks = 0;

-- ----------------------------
-- table structure for messages
-- ----------------------------
drop table if exists `messages`;
create table `messages`
(
    `id`         int(11) unsigned                                              not null auto_increment,
    `user_id`    int(11)                                                       not null comment '用户id',
    `room_id`    int(11)                                                       not null comment '房间id',
    `to_user_id` int(11)                                                       null default 0 comment '私聊用户id',
    `content`    longtext character set utf8mb4 collate utf8mb4_general_ci     null comment '聊天内容',
    `image_url`  varchar(255) character set utf8mb4 collate utf8mb4_general_ci null default '' comment '图片url',
    `created_at` datetime(0)                                                   null default null,
    `updated_at` datetime(0)                                                   null default null on update current_timestamp(0),
    `deleted_at` datetime(0)                                                   null default null,
    primary key (`id`) using btree,
    index `idx_user_id` (`user_id`) using btree
) engine = innodb
  character set = utf8mb4
  collate = utf8mb4_general_ci
  row_format = dynamic;

-- ----------------------------
-- table structure for users
-- ----------------------------
drop table if exists `users`;
create table `users`
(
    `id`         int(11) unsigned                                              not null auto_increment,
    `username`   varchar(50) character set utf8mb4 collate utf8mb4_general_ci  not null default '' comment '昵称',
    `password`   varchar(125) character set utf8mb4 collate utf8mb4_general_ci null     default '' comment '密码',
    `avatar_id`  varchar(50) character set utf8mb4 collate utf8mb4_general_ci  null     default '1' comment '头像id',
    `created_at` datetime(0)                                                   null     default null,
    `updated_at` datetime(0)                                                   null     default null on update current_timestamp(0),
    `deleted_at` datetime(0)                                                   null     default null,
    primary key (`id`) using btree,
    unique index `username` (`username`) using btree
) engine = innodb
  character set = utf8mb4
  collate = utf8mb4_general_ci
  row_format = dynamic;

set foreign_key_checks = 1;
