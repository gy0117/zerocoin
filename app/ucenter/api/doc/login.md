## ucenter-api服务

- DEV：0.0.0.0:8888

## /uc/login

### 登录

#### 请求方法: POST

#### 请求参数

```json
{
  "username": "xingzhi9",
  "password": "123456"
}
```

#### 响应参数

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "username": "xingzhi9",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA5MDE4NDgsImlhdCI6MTcxMDI5NzA0OCwidXNlcklkIjoxOX0.OiYsqIp5i7AUmmf00ZrqITp-8maeJfE0fntndDjGqic",
        "memberLevel": "普通会员",
        "realName": "",
        "country": "",
        "avatar": "https://p2.itc.cn/q_70/images03/20230902/721191166cd2465c9db74d5b52a3e7bc.png",
        "promotionCode": "",
        "id": 19,
        "loginCount": 0,
        "superPartner": "0",
        "memberRate": 0
    }
}
```


## /uc/check/login

### 检查登录状态

#### 请求方法: POST

#### 请求参数

header中带上token

| X-Auth-Token | xxx  |
| ------------ | ---- |



#### 响应参数

```json
{
    "code": 0,
    "message": "success",
    "data": true
}
```