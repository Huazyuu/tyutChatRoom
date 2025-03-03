# 项目 II：即时通信（IM）应用开发(只写了后端)

>就是一个练手项目,没有答辩,我们组做的项目一,我自己写写后端练练手

>另一个版本[another version](AnotherVersion/readme.md)

---

## require

前后端分离,只有后端

go version >= 1.20

安装需要依赖 `go mod tidy`

修改`conf/settings.yaml`配置文件中的 mysql redis log system email(不想使用邮箱验证可以去代码里删掉)

---



## 进度

### 用户

#### 注册

`http://127.0.0.1:8080/api/users/register`

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

`http://127.0.0.1:8080/api/users/login`

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

`http://127.0.0.1:8080/api/users/logout`

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

`http://127.0.0.1:8080/api/users`

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
params Authorization:bearer xxx.xxx.xxx

```js
    const jwtToken = 'bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Iua1i-ivleeUqOaItyIsInVzZXJfaWQiOiI5N2FiMDk2ZTJhIiwiZXhwIjoxNzQyNDgwNTk5LjI5NDA0MjgsImlzcyI6Inp5dXUifQ.iMYaKVNcVrS-lejRgZYGUU33CihquKhivPB5M3CwBNA';
    const baseWsServerUrl = 'ws://127.0.0.1:8080/api/chat_groups';
    const wsServerUrl = `${baseWsServerUrl}?Authorization=${jwtToken}`;

    // 建立 WebSocket 连接
    const socket = new WebSocket(wsServerUrl);

    // 监听 WebSocket 连接打开事件
    socket.addEventListener('open', (event) => {
        console.log('WebSocket 连接已打开');
        const data = {
            "content": "hello google",
            "msg_type": 4
        };

        // 将 JavaScript 对象转换为 JSON 字符串
        const jsonData = JSON.stringify(data);

        // 发送 JSON 数据
        socket.send(jsonData);

        console.log('数据已发送');
    });

```

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
params Authorization:bearer xxx.xxx.xxx

`ws://127.0.0.1:8080/api/chat_private`

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


#### 聊天记录

需要用户携带token 在http request header
Authorization:bearer xxx.xxx.xxx

##### 群聊

`http://127.0.0.1:8080/api/chat_groupList`


| params   | value                                                    |
| -------- | -------------------------------------------------------- |
| limit    | 10                                                       |
| page     | 1                                                        |
| sort     | username asc/desc                                        |
| username | 输入用户名查找其群聊记录,用户名为空或错误,将显示所有记录 |

##### 私聊

`http://127.0.0.1:8080/api/chat_privateList`


| params   | value                          |
| -------- | ------------------------------ |
| limit    | 10                             |
| page     | 1                              |
| sort     | username asc/desc              |
| username | 输入用户名查找于其相关私聊记录 |

---



### 文件管理

#### 上传

需要BearerToken

| params    | value         |
| --------- | ------------- |
| target_id | 10位的user_id |

| form-data | value    |
| --------- | -------- |
| file      | 文件对象 |

**resp**

```json
{
    "code": 0,
    "data": {
        "user_id": "dfe0b20eb9",
        "target_id": "97ab096e2a",
        "path": "uploads\\file\\adminUser\\1740377182215261100_《计算机网络课程设计》课设目标.md",
        "file_name": "1740377182215261100_《计算机网络课程设计》课设目标.md",
        "file_size": 6915,
        "file_type": "md"
    },
    "msg": "成功"
}
```





#### 下载

需要BearerToken

| params | value  |
| ------ | ------ |
| file   | 文件名 |

返回文件对象

#### ws调用文件相关http

```
// 上传
//ws req
{
    "content": "hello test",
    "msg_type": 5,
    "file": {
        "path": "C:/Users/yu/Desktop/设计模式.txt",
        "name": "ws上传下载测试.txt",
        "type": "md",
        "size": 0
    }
}

// ws resp
{
    "user_id": "97ab096e2a",
    "username": "测试用户",
    "avatar": "uploads/avatar/测试用户.png",
    "target_id": "dfe0b20eb9",
    "target_name": "adminUser",
    "target_avatar": "uploads/avatar/adminUser.png",
    "msg_type": 5,
    "content": "C:/Users/yu/Desktop/设计模式.txt",
    "date": "2025-02-24T17:19:29.5441526+08:00",
    "online_count": 2
}
```

```下载
// 下载
//req 
{
    "content": "hello test",
    "msg_type": 6,
    "file": {
        "name": "ws上传下载测试2.txt",
        "type": "md",
        "size": 0
    }
}
// resp 返回文件对象
// resp info json
{
    "user_id": "dfe0b20eb9",
    "username": "adminUser",
    "avatar": "uploads/avatar/adminUser.png",
    "target_id": "97ab096e2a",
    "target_name": "测试用户",
    "target_avatar": "uploads/avatar/测试用户.png",
    "msg_type": 6,
    "content": "",
    "date": "2025-02-24T17:25:26.8921766+08:00",
    "online_count": 2
}
```

