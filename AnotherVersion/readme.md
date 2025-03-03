# go-gin-chat(Gin+websocket(gorilla) 的多人聊天室后端)

> 借鉴于 [github地址](https://github.com/hezhizheng/go-gin-chat)

## Feature
- 登录/注册(防止重复登录)
- 群聊(多房间、支持文字、emoji、文件(图片)上传，使用 [freeimage.host](https://freeimage.host/)做图床 )
- 私聊(消息提醒)
- 历史消息查看(点击加载更多)
- 心跳检测及自动重连机制，来自 https://github.com/zimv/websocket-heartbeat-js
- go mod 包管理
- 前后端分离
- 使用zap日志记录
- 支持 http/ws 、 https/wss

## 结构

```
├─conf
│      conf.go
│      settings.yaml
│      
├─controller
│      chat.go
│      home.go
│      img.go
│      index.go
│      page.go
│      room.go
│      user.go
│      ws.go
│      
├─doc
│      freeimg.md
│      
├─log
│      log.go
│      
├─logdoc
│      loginfo.log
│      
├─middleware
│  └─session
│          session.go
│          
├─models
│  │  message.go
│  │  mysql.go
│  │  user.go
│  │  
│  ├─req
│  │      user.go
│  │      
│  └─res
│          response.go
│          
├─routes
│      route.go
│      
├─service
│  ├─img_service
│  │      freeimagehost.go
│  │      freeimagehost_test.go
│  │      img.go
│  │      smmsapp.go
│  │      
│  ├─message_service
│  │      msg.go
│  │      
│  └─user_service
│          user.go
│          
├─test
│      imgfreeLocal.go
│      imgfreeUrl.go
│      
├─uploads
│      QQ图片20230103083100.jpg
│      
├─utils
│      in_array.go
│      pwd_encrypt.go
│      substr_len.go
│      threadSafety.go
│      valid_msg.go
│      
└─ws
        serve.go
        ServeInterface.go
        variable.go
```



Base URLs:

# Authentication

* API Key (apikey-header-token)
  - Parameter Name: **token**, in: header. 

# ginchat

## POST 注册登录

POST /login

> Body 请求参数

```json
{
  "username": "string",
  "password": "string",
  "avatar_id": "string"
}
```

### 请求参数

| 名称        | 位置 | 类型   | 必选 | 说明 |
| ----------- | ---- | ------ | ---- | ---- |
| body        | body | object | 否   | none |
| » username  | body | string | 是   | none |
| » password  | body | string | 是   | none |
| » avatar_id | body | string | 是   | none |

> 返回示例

```json
{
  "code": 0,
  "data": {},
  "msg": "登陆成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称   | 类型    | 必选 | 约束 | 中文名 | 说明 |
| ------ | ------- | ---- | ---- | ------ | ---- |
| » code | integer | true | none |        | none |
| » data | object  | true | none |        | none |
| » msg  | string  | true | none |        | none |

## GET 退出登录

GET /logout

> 返回示例

```json
{
  "code": 0,
  "data": {},
  "msg": "退出登录"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称   | 类型    | 必选 | 约束 | 中文名 | 说明 |
| ------ | ------- | ---- | ---- | ------ | ---- |
| » code | integer | true | none |        | none |
| » data | object  | true | none |        | none |
| » msg  | string  | true | none |        | none |

## GET index

GET /

> 返回示例

```json
{
  "code": 0,
  "data": {
    "rooms": [
      {
        "id": 1,
        "num": 0
      },
      {
        "id": 2,
        "num": 0
      },
      {
        "id": 3,
        "num": 0
      },
      {
        "id": 4,
        "num": 0
      },
      {
        "id": 5,
        "num": 0
      },
      {
        "id": 6,
        "num": 0
      }
    ],
    "user_info": {
      "avatar_id": "1231",
      "uid": 4,
      "username": "user1"
    }
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称          | 类型     | 必选 | 约束 | 中文名 | 说明 |
| ------------- | -------- | ---- | ---- | ------ | ---- |
| » code        | integer  | true | none |        | none |
| » data        | object   | true | none |        | none |
| »» rooms      | [object] | true | none |        | none |
| »»» id        | integer  | true | none |        | none |
| »»» num       | integer  | true | none |        | none |
| »» user_info  | object   | true | none |        | none |
| »»» avatar_id | string   | true | none |        | none |
| »»» uid       | integer  | true | none |        | none |
| »»» username  | string   | true | none |        | none |
| » msg         | string   | true | none |        | none |

## GET 历史记录

GET /pagination

### 请求参数

| 名称    | 位置  | 类型   | 必选 | 说明 |
| ------- | ----- | ------ | ---- | ---- |
| room_id | query | string | 否   | none |
| uid     | query | string | 否   | none |
| offset  | query | string | 否   | none |

> 返回示例

```json
{
  "code": 0,
  "data": {
    "count": 100,
    "list": [
      {
        "avatar_id": "1",
        "content": "test user 2",
        "created_at": "2025-03-03T19:37:51+08:00",
        "deleted_at": null,
        "id": 134,
        "image_url": "",
        "room_id": 1,
        "to_user_id": 4,
        "updated_at": "2025-03-03T19:37:51+08:00",
        "user_id": 5,
        "username": "user2"
      }
    ]
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称           | 类型     | 必选  | 约束 | 中文名 | 说明 |
| -------------- | -------- | ----- | ---- | ------ | ---- |
| » code         | integer  | true  | none |        | none |
| » data         | object   | true  | none |        | none |
| »» count       | integer  | true  | none |        | none |
| »» list        | [object] | true  | none |        | none |
| »»» avatar_id  | string   | false | none |        | none |
| »»» content    | string   | false | none |        | none |
| »»» created_at | string   | false | none |        | none |
| »»» deleted_at | null     | false | none |        | none |
| »»» id         | integer  | false | none |        | none |
| »»» image_url  | string   | false | none |        | none |
| »»» room_id    | integer  | false | none |        | none |
| »»» to_user_id | integer  | false | none |        | none |
| »»» updated_at | string   | false | none |        | none |
| »»» user_id    | integer  | false | none |        | none |
| »»» username   | string   | false | none |        | none |
| » msg          | string   | true  | none |        | none |

## GET 主界面

GET /home

> 返回示例

```json
{
  "code": 0,
  "data": {
    "rooms": [
      {
        "id": 1,
        "num": 0
      },
      {
        "id": 2,
        "num": 0
      },
      {
        "id": 3,
        "num": 0
      },
      {
        "id": 4,
        "num": 0
      },
      {
        "id": 5,
        "num": 0
      },
      {
        "id": 6,
        "num": 0
      }
    ],
    "user_info": {
      "avatar_id": "1231",
      "uid": 1,
      "username": "盛梓馨"
    }
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称          | 类型     | 必选  | 约束 | 中文名 | 说明 |
| ------------- | -------- | ----- | ---- | ------ | ---- |
| » code        | integer  | true  | none |        | none |
| » data        | object   | true  | none |        | none |
| »» rooms      | [object] | true  | none |        | none |
| »»» id        | integer  | true  | none |        | none |
| »»» num       | integer  | true  | none |        | none |
| »» user_info  | object   | true  | none |        | none |
| »»» avatar_id | string   | true  | none |        | none |
| »»» uid       | integer  | true  | none |        | none |
| »»» username  | string   | true  | none |        | none |
| » msg         | string   | true  | none |        | none |
| » *anonymous* | string   | false | none |        | none |

## GET 房间

GET /room/{roomid}

### 请求参数

| 名称   | 位置 | 类型   | 必选 | 说明 |
| ------ | ---- | ------ | ---- | ---- |
| roomid | path | string | 是   | none |

> 返回示例

```json
{
  "code": 0,
  "data": {
    "msg_list": [
      {
        "avatar_id": "1",
        "content": "user test1 private",
        "created_at": "2025-03-03T18:03:04+08:00",
        "deleted_at": null,
        "id": 67,
        "image_url": "",
        "room_id": 1,
        "to_user_id": 0,
        "updated_at": "2025-03-03T18:03:04+08:00",
        "user_id": 4,
        "username": "user1"
      }
    ],
    "msg_list_count": 40,
    "room_id": "1",
    "user_info": {
      "avatar_id": "1231",
      "uid": 1,
      "username": "盛梓馨"
    }
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称              | 类型     | 必选 | 约束 | 中文名 | 说明 |
| ----------------- | -------- | ---- | ---- | ------ | ---- |
| » code            | integer  | true | none |        | none |
| » data            | object   | true | none |        | none |
| »» msg_list       | [object] | true | none |        | none |
| »»» avatar_id     | string   | true | none |        | none |
| »»» content       | string   | true | none |        | none |
| »»» created_at    | string   | true | none |        | none |
| »»» deleted_at    | null     | true | none |        | none |
| »»» id            | integer  | true | none |        | none |
| »»» image_url     | string   | true | none |        | none |
| »»» room_id       | integer  | true | none |        | none |
| »»» to_user_id    | integer  | true | none |        | none |
| »»» updated_at    | string   | true | none |        | none |
| »»» user_id       | integer  | true | none |        | none |
| »»» username      | string   | true | none |        | none |
| »» msg_list_count | integer  | true | none |        | none |
| »» room_id        | string   | true | none |        | none |
| »» user_info      | object   | true | none |        | none |
| »»» avatar_id     | string   | true | none |        | none |
| »»» uid           | integer  | true | none |        | none |
| »»» username      | string   | true | none |        | none |
| » msg             | string   | true | none |        | none |

## GET 私聊记录

GET /private-chat

> 返回示例

```json
{
  "code": 0,
  "data": {
    "msg_list": [
      {
        "avatar_id": "1231",
        "content": "user test1 private",
        "created_at": "2025-03-03T18:03:04+08:00",
        "deleted_at": null,
        "id": 67,
        "image_url": "",
        "room_id": 1,
        "to_user_id": 0,
        "updated_at": "2025-03-03T18:03:04+08:00",
        "user_id": 4,
        "username": "user1"
      }
    ],
    "room_id": "",
    "user_info": {
      "avatar_id": "1231",
      "uid": 4,
      "username": "user1"
    }
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称           | 类型     | 必选  | 约束 | 中文名 | 说明 |
| -------------- | -------- | ----- | ---- | ------ | ---- |
| » code         | integer  | true  | none |        | none |
| » data         | object   | true  | none |        | none |
| »» msg_list    | [object] | true  | none |        | none |
| »»» avatar_id  | string   | false | none |        | none |
| »»» content    | string   | false | none |        | none |
| »»» created_at | string   | false | none |        | none |
| »»» deleted_at | null     | false | none |        | none |
| »»» id         | integer  | false | none |        | none |
| »»» image_url  | string   | false | none |        | none |
| »»» room_id    | integer  | false | none |        | none |
| »»» to_user_id | integer  | false | none |        | none |
| »»» updated_at | string   | false | none |        | none |
| »»» user_id    | integer  | false | none |        | none |
| »»» username   | string   | false | none |        | none |
| »» room_id     | string   | true  | none |        | none |
| »» user_info   | object   | true  | none |        | none |
| »»» avatar_id  | string   | true  | none |        | none |
| »»» uid        | integer  | true  | none |        | none |
| »»» username   | string   | true  | none |        | none |
| » msg          | string   | true  | none |        | none |

## POST 图片上传到图床

POST /img-upload

> Body 请求参数

```yaml
file: file://E:\OneDrive - exsn4848\图片\QQ图片20230103083100.jpg

```

### 请求参数

| 名称   | 位置 | 类型           | 必选 | 说明 |
| ------ | ---- | -------------- | ---- | ---- |
| body   | body | object         | 否   | none |
| » file | body | string(binary) | 是   | none |

> 返回示例

```json
{
  "code": 0,
  "data": {
    "code": 0,
    "data": {
      "url": "https://iili.io/33o3vOg.md.jpg"
    }
  },
  "msg": "成功"
}
```

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

### 返回数据结构

状态码 **200**

| 名称    | 类型    | 必选 | 约束 | 中文名 | 说明 |
| ------- | ------- | ---- | ---- | ------ | ---- |
| » code  | integer | true | none |        | none |
| » data  | object  | true | none |        | none |
| »» code | integer | true | none |        | none |
| »» data | object  | true | none |        | none |
| »»» url | string  | true | none |        | none |
| » msg   | string  | true | none |        | none |

# 

