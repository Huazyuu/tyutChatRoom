# 《计算机网络课程设计》课设目标

《计算机网络课程设计》是软件工程专业的实践课程之一，是在学习完《计算机网络》课程后开展的全面综合练习。通过规划和组建简单网络，学生将达成以下知识与能力目标：

---

## 知识目标

1. **理解网络基本概念与结构**：理解计算机网络的基本概念，掌握网络结构以及常见网络设备的功能，强化对计算机网络软、硬件组成的感性认识。
2. **掌握网络协议与标准**：学会使用网络协议和标准，了解 TCP/IP 协议族的原理和应用。深入理解网络分层结构概念，特别是会话层、表示层、应用层等高层协议软件的通信功能和实现方法。掌握网络互连设备的使用及工作原理，以及 IP 地址的配置。
3. **掌握故障排除与配置技巧**：掌握网络故障排除方法和网络设备配置技巧。

## 能力目标

1. **设计与搭建小型网络**：能够运用所学知识设计和搭建小型网络，实现设备互联和数据通信。
2. **熟练使用网络诊断工具**：熟练使用网络诊断工具，分析网络故障原因并解决问题。
3. **掌握局域网设计技术**：掌握局域网的设计技术和技巧，培养开发网络应用的独立工作能力。
4. **提高团队协作能力**：通过分组合作完成网络实践项目，提高团队协作能力。

# 项目 II：即时通信（IM）应用开发

## 2.2.1 备选题目及需求

设计一个基于 TCP/IP 协议的网络通讯小应用，可采用 UDP 或 TCP 实现，功能需满足以下要求：

- **架构要求**：必须是 C/S 应用，编程语言、开发库、开发框架等不限。
- **用户管理**：能够进行用户管理，所有用户必须登录到服务器，由服务器维护在线信息。
- **实时通信**：用户登录后能够进行实时的多方点到点短信息通信，如聊天。
- **聊天记录**：能够保存聊天记录到数据库，并允许用户查看历史聊天记录。
- **文件传输**：允许用户间进行文件传输，传输过程中实时显示完成进度。



## 进度

### 用户

#### 注册

```
// request 
{
    "email": "testemail@outlook.com",
    "username": "测试用户",
    "password": "zyuuforyu"
}

// response
{
    "code": 0,
    "data": {},
    "msg": "用户创建成功 测试用户 "
}
```

#### 登录

```
// 第一次request 
{
    "email": "testemail@outlook.com",
    "password": "zyuuforyu"
}
// 第一次resp
{
    "code": 0,
    "data": {},
    "msg": "验证码已发送"
}
```

第二次请求会收到验证码 

>需要自行配置邮件发送插件
>
>参考  [gomail库](https://pkg.go.dev/gopkg.in/gomail.v2) 教程配置

```
// req
{
    "email": "zyuuforyu@outlook.com",
    "password": "zyuuforyu",
    "code": "0661"
}

// resp
{
    "code": 0,
    "data": {},
    "msg": "用户创建成功 测试用户 "
}
```

#### 注销

不需要 req json

前端需要再request header中添加token

建议前端持久化在Cookie中 当然后端数据库中也存储了token

```
// resp
{
    "code": 0,
    "data": {},
    "msg": "注销成功"
}
```
#### 用户列表

| params | value             |
| ------ | ----------------- |
| limit  | 10                |
| page   | 1                 |
| sort   | username asc/desc |

```
// resp
{
    "code": 0,
    "data": {
        "count": 1,
        "list": [
            {
                "created_at": "0001-01-01T00:00:00Z",
                "updated_at": "0001-01-01T00:00:00Z",
                "user_id": "2a67f4734d",
                "user_name": "第二个人",
                "avatar": "uploads/avatar/第二个人.png",
                "email": "fmnnzu512@outlook.com"
            },
        ]
    },
    "msg": "成功"
}

```


### 聊天室 websocket

#### 群聊

需要用户携带token

`ws://127.0.0.1:8080/api/chat_groups`

| 类型(int)       | 值   |
| --------------- | ---- |
| SystemMsg       | 1    |
| InRoomMsg2      | 2    |
| OutRoomMsg      | 3    |
| TextMsg         | 4    |
| FileMsg         | 5    |
| FileProgressMsg | 6    |
| ImageMsg        | 7    |



```
// req
{
    "content": "hello 2号",
    "msg_type": 4
}

// resp
{
    "user_id": "97ab096e2a",
    "username": "测试用户",
    "avatar": "uploads/avatar/测试用户.png",
    "msg_type": 4,
    "content": "hello 2号",
    "date": "2025-02-19T16:29:52.6806232+08:00",
    "online_count": 1
}
```
#### 私聊

需要用户携带token

```
ws://127.0.0.1:8080/api/chat_private?target_id=xxxx
需要传递param query target_id
```

| 类型(int)       | 值   |
| --------------- | ---- |
| SystemMsg       | 1    |
| InRoomMsg2      | 2    |
| OutRoomMsg      | 3    |
| TextMsg         | 4    |
| FileMsg         | 5    |
| FileProgressMsg | 6    |
| ImageMsg        | 7    |

```
// req
{
    "content": "hello 2号",
    "msg_type": 4
}

// resp
{
    "user_id": "97ab096e2a",
    "username": "测试用户",
    "avatar": "uploads/avatar/测试用户.png",
    "target_id": "dfe0b20eb9",
    "target_name": "adminUser",
    "target_avatar": "uploads/avatar/adminUser.png",
    "msg_type": 4,
    "content": "hello 2号",
    "date": "2025-02-19T22:10:58.5416886+08:00",
    "online_count": 2
}
```


