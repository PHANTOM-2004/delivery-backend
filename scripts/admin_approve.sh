#!/bin/sh

set -x

# 创建一个管理员
curl "http://localhost:8000/api/v1/admin/create?super_token=LTDZ&admin_name=cjq&password=123451234512345&account=abcabcabcabc" --request POST

# 管理员登录
curl -c cookies.txt -F "account=abcabcabcabc" -F "password=123451234512345" http://localhost:8000/api/v1/admin/login

# 验证登录状态
curl -b cookies.txt http://localhost:8000/api/v1/admin/login-status

# 上传申请表
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=szc" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST

# 上传申请表
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=mhc" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST

# 上传申请表
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=xp" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST

# 上传申请表
curl \
  -F "description=ok" \
  -F "email=666@qq.com" \
  -F "phone_number=+8618537775175" \
  -F "license=@Makefile.png" \
  -F "name=ych" \
  http://localhost:8000/api/v1/customer/merchant-application -H "Content-Type: multipart/form-data" \
  --request POST

# 获得申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1

# 批准申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/1/approve --request PUT

# 批准申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/2/approve --request PUT

# 批准申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/3/approve --request PUT

# 批准申请表
curl -b cookies.txt http://localhost:8000/api/v1/admin/merchant-application/4/approve --request PUT
