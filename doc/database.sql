create database ginGorilla_db;
use ginGorilla_db;

-- 创建用户表
create table if not exists user_models
(
    id         int auto_increment primary key,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,


    user_id    varchar(10)  not null unique,
    username   varchar(50)  not null,
    password   varchar(255) not null,
    avatar     varchar(255),
    email      varchar(100) unique,
    token      varchar(255)
);

-- 创建聊天记录表
create table if not exists chat_models
(
    id         int auto_increment primary key,
    created_at timestamp            default current_timestamp,
    updated_at timestamp            default current_timestamp on update current_timestamp,

    user_id    varchar(20) not null,
    target_id  varchar(20) not null,
    content    text        not null,

    ip         varchar(45),
    addr       varchar(255),

    is_group   boolean     not null default false,
    msg_type   int         not null,

    foreign key (user_id) references user_models (user_id),
    foreign key (target_id) references user_models (user_id)
);
SELECT chat_models.user_id,
       chat_models.content,
       user.username,
       chat_models.created_at,
       chat_models.ip,
       chat_models.addr
FROM `chat_models`
         LEFT JOIN user_models AS user ON user.user_id = chat_models.user_id
WHERE chat_models.is_group = 1
LIMIT 10;
