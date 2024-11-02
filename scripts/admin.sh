#!/bin/bash

# 创建一个管理员
curl "https://localhost/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=0123456789" --request POST

# 管理员登录
curl -c cookies.txt -F "account=0123456789" -F "password=123451234512345" https://localhost/admin/login

# 鉴权
curl -b cookies.txt https://localhost/api/v1/admin/jwt/auth --request GET

# 获得申请表
curl -b cookies.txt https://localhost/api/v1/admin/jwt/merchant-application/1

# 上传申请表
►curl -F "status=1" \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  "https://localhost/api/v1/customer/business-application" \
  -H "Content-Type: multipart/form-data" \
  --request POST
