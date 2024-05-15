## ucenter-api服务
- DEV：0.0.0.0:8888

## /uc/register/phone

### 注册

#### 请求方法: POST
#### 请求参数
```json
{
  "phone": "15996230001",
  "username": "xingzhi",
  "password": "123456"
}
```

#### 响应参数
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "accessToken": "",
        "accessExpire": 0,
        "refreshAfter": 0
    }
}
```


## /uc/mobile/code

### 发送验证码
#### 请求方法: POST
#### 请求参数
```json
{
    "country": "+86",
    "phone": "15910088378"
}
```

#### 响应参数
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "smsCode": "1228"
    }
}
```