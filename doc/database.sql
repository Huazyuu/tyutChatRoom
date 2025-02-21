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

-- 创建文件表
-- 创建文件表
create table if not exists file_models (
    id int auto_increment primary key,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    
    user_id varchar(255) not null,
    target_id varchar(255) not null,
    path varchar(255) not null,
    file_name varchar(255) not null,
    file_size bigint not null,
    file_type varchar(255) not null,
    
    -- 定义外键约束
    foreign key (user_id) references user_models(user_id) on delete cascade on update cascade,
    foreign key (target_id) references user_models(user_id) on delete cascade on update cascade
);